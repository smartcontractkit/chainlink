// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vrf_coordinator_v2_events_mock

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

type VRFCoordinatorV2EventsMockFeeConfig struct {
	FulfillmentFlatFeeLinkPPMTier1 uint32
	FulfillmentFlatFeeLinkPPMTier2 uint32
	FulfillmentFlatFeeLinkPPMTier3 uint32
	FulfillmentFlatFeeLinkPPMTier4 uint32
	FulfillmentFlatFeeLinkPPMTier5 uint32
	ReqsForTier2                   *big.Int
	ReqsForTier3                   *big.Int
	ReqsForTier4                   *big.Int
	ReqsForTier5                   *big.Int
}

var VRFCoordinatorV2EventsMockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"indexed\":false,\"internalType\":\"structVRFCoordinatorV2EventsMock.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FundsRecovered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyDeregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"ProvingKeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"RandomWordsFulfilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"SubscriptionCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"SubscriptionConsumerRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"SubscriptionCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"SubscriptionFunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"SubscriptionOwnerTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"maxGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"stalenessSeconds\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"gasAfterPaymentCalculation\",\"type\":\"uint32\"},{\"internalType\":\"int256\",\"name\":\"fallbackWeiPerUnitLink\",\"type\":\"int256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier1\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier2\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier3\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier4\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"fulfillmentFlatFeeLinkPPMTier5\",\"type\":\"uint32\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier2\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier3\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier4\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"reqsForTier5\",\"type\":\"uint24\"}],\"internalType\":\"structVRFCoordinatorV2EventsMock.FeeConfig\",\"name\":\"feeConfig\",\"type\":\"tuple\"}],\"name\":\"emitConfigSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emitFundsRecovered\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"emitProvingKeyDeregistered\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"emitProvingKeyRegistered\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"outputSeed\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"payment\",\"type\":\"uint96\"},{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"emitRandomWordsFulfilled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"preSeed\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"minimumRequestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"emitRandomWordsRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emitSubscriptionCanceled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"emitSubscriptionConsumerAdded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"consumer\",\"type\":\"address\"}],\"name\":\"emitSubscriptionConsumerRemoved\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"emitSubscriptionCreated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"oldBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newBalance\",\"type\":\"uint256\"}],\"name\":\"emitSubscriptionFunded\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitSubscriptionOwnerTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"subId\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitSubscriptionOwnerTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610c4d806100206000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c8063b019b4e811610097578063e144e45a11610066578063e144e45a146101cd578063e2cab57b146101e0578063f7420bc2146101f3578063fe62d3e91461020657600080fd5b8063b019b4e814610181578063ca920adb14610194578063dde69b3f146101a7578063e0f6eff1146101ba57600080fd5b80634bf6a80d116100d35780634bf6a80d14610135578063675b924414610148578063689300ea1461015b578063a52572261461016e57600080fd5b80631917c3ed146100fa5780633f70afb61461010f578063438cbfbb14610122575b600080fd5b61010d61010836600461080e565b610219565b005b61010d61011d366004610a50565b61026d565b61010d61013036600461080e565b6102bf565b61010d610143366004610a6c565b610307565b61010d610156366004610a50565b610369565b61010d6101693660046107e4565b6103bb565b61010d61017c366004610a50565b61040e565b61010d61018f3660046107b1565b610460565b61010d6101a2366004610831565b6104be565b61010d6101b5366004610aaf565b610546565b61010d6101c8366004610a6c565b61059e565b61010d6101db3660046108b3565b6105f7565b61010d6101ee366004610aeb565b610640565b61010d6102013660046107b1565b610682565b61010d6102143660046109f0565b6106e0565b8073ffffffffffffffffffffffffffffffffffffffff167f72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d8360405161026191815260200190565b60405180910390a25050565b60405173ffffffffffffffffffffffffffffffffffffffff8216815267ffffffffffffffff8316907f464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf90602001610261565b8073ffffffffffffffffffffffffffffffffffffffff167fe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b88360405161026191815260200190565b6040805173ffffffffffffffffffffffffffffffffffffffff80851682528316602082015267ffffffffffffffff8516917f6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f091015b60405180910390a2505050565b60405173ffffffffffffffffffffffffffffffffffffffff8216815267ffffffffffffffff8316907f43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e090602001610261565b6040805173ffffffffffffffffffffffffffffffffffffffff84168152602081018390527f59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600910160405180910390a15050565b60405173ffffffffffffffffffffffffffffffffffffffff8216815267ffffffffffffffff8316907f182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b90602001610261565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b604080518881526020810188905261ffff86168183015263ffffffff858116606083015284166080820152905173ffffffffffffffffffffffffffffffffffffffff83169167ffffffffffffffff8816918b917f63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772919081900360a00190a45050505050505050565b6040805173ffffffffffffffffffffffffffffffffffffffff841681526020810183905267ffffffffffffffff8516917fe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815910161035c565b6040805173ffffffffffffffffffffffffffffffffffffffff80851682528316602082015267ffffffffffffffff8516917f69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be910161035c565b7fc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb286868686868660405161063096959493929190610b1e565b60405180910390a1505050505050565b604080518381526020810183905267ffffffffffffffff8516917fd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8910161035c565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b604080518481526bffffffffffffffffffffffff8416602082015282151581830152905185917f7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4919081900360600190a250505050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461075b57600080fd5b919050565b803561ffff8116811461075b57600080fd5b803562ffffff8116811461075b57600080fd5b803563ffffffff8116811461075b57600080fd5b803567ffffffffffffffff8116811461075b57600080fd5b600080604083850312156107c457600080fd5b6107cd83610737565b91506107db60208401610737565b90509250929050565b600080604083850312156107f757600080fd5b61080083610737565b946020939093013593505050565b6000806040838503121561082157600080fd5b823591506107db60208401610737565b600080600080600080600080610100898b03121561084e57600080fd5b88359750602089013596506040890135955061086c60608a01610799565b945061087a60808a01610760565b935061088860a08a01610785565b925061089660c08a01610785565b91506108a460e08a01610737565b90509295985092959890939650565b6000806000806000808688036101c08112156108ce57600080fd5b6108d788610760565b96506108e560208901610785565b95506108f360408901610785565b945061090160608901610785565b935060808801359250610120807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff608301121561093c57600080fd5b610944610bef565b915061095260a08a01610785565b825261096060c08a01610785565b602083015261097160e08a01610785565b6040830152610100610984818b01610785565b6060840152610994828b01610785565b60808401526109a66101408b01610772565b60a08401526109b86101608b01610772565b60c08401526109ca6101808b01610772565b60e08401526109dc6101a08b01610772565b818401525050809150509295509295509295565b60008060008060808587031215610a0657600080fd5b843593506020850135925060408501356bffffffffffffffffffffffff81168114610a3057600080fd5b915060608501358015158114610a4557600080fd5b939692955090935050565b60008060408385031215610a6357600080fd5b6107cd83610799565b600080600060608486031215610a8157600080fd5b610a8a84610799565b9250610a9860208501610737565b9150610aa660408501610737565b90509250925092565b600080600060608486031215610ac457600080fd5b610acd84610799565b9250610adb60208501610737565b9150604084013590509250925092565b600080600060608486031215610b0057600080fd5b610b0984610799565b95602085013595506040909401359392505050565b60006101c08201905061ffff8816825263ffffffff8088166020840152808716604084015280861660608401528460808401528084511660a08401528060208501511660c0840152506040830151610b7e60e084018263ffffffff169052565b506060830151610100610b988185018363ffffffff169052565b608085015163ffffffff1661012085015260a085015162ffffff90811661014086015260c0860151811661016086015260e086015181166101808601529401519093166101a0909201919091529695505050505050565b604051610120810167ffffffffffffffff81118282101715610c3a577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040529056fea164736f6c6343000806000a",
}

var VRFCoordinatorV2EventsMockABI = VRFCoordinatorV2EventsMockMetaData.ABI

var VRFCoordinatorV2EventsMockBin = VRFCoordinatorV2EventsMockMetaData.Bin

func DeployVRFCoordinatorV2EventsMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFCoordinatorV2EventsMock, error) {
	parsed, err := VRFCoordinatorV2EventsMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV2EventsMockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorV2EventsMock{VRFCoordinatorV2EventsMockCaller: VRFCoordinatorV2EventsMockCaller{contract: contract}, VRFCoordinatorV2EventsMockTransactor: VRFCoordinatorV2EventsMockTransactor{contract: contract}, VRFCoordinatorV2EventsMockFilterer: VRFCoordinatorV2EventsMockFilterer{contract: contract}}, nil
}

type VRFCoordinatorV2EventsMock struct {
	address common.Address
	abi     abi.ABI
	VRFCoordinatorV2EventsMockCaller
	VRFCoordinatorV2EventsMockTransactor
	VRFCoordinatorV2EventsMockFilterer
}

type VRFCoordinatorV2EventsMockCaller struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2EventsMockTransactor struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2EventsMockFilterer struct {
	contract *bind.BoundContract
}

type VRFCoordinatorV2EventsMockSession struct {
	Contract     *VRFCoordinatorV2EventsMock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV2EventsMockCallerSession struct {
	Contract *VRFCoordinatorV2EventsMockCaller
	CallOpts bind.CallOpts
}

type VRFCoordinatorV2EventsMockTransactorSession struct {
	Contract     *VRFCoordinatorV2EventsMockTransactor
	TransactOpts bind.TransactOpts
}

type VRFCoordinatorV2EventsMockRaw struct {
	Contract *VRFCoordinatorV2EventsMock
}

type VRFCoordinatorV2EventsMockCallerRaw struct {
	Contract *VRFCoordinatorV2EventsMockCaller
}

type VRFCoordinatorV2EventsMockTransactorRaw struct {
	Contract *VRFCoordinatorV2EventsMockTransactor
}

func NewVRFCoordinatorV2EventsMock(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorV2EventsMock, error) {
	abi, err := abi.JSON(strings.NewReader(VRFCoordinatorV2EventsMockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindVRFCoordinatorV2EventsMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMock{address: address, abi: abi, VRFCoordinatorV2EventsMockCaller: VRFCoordinatorV2EventsMockCaller{contract: contract}, VRFCoordinatorV2EventsMockTransactor: VRFCoordinatorV2EventsMockTransactor{contract: contract}, VRFCoordinatorV2EventsMockFilterer: VRFCoordinatorV2EventsMockFilterer{contract: contract}}, nil
}

func NewVRFCoordinatorV2EventsMockCaller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorV2EventsMockCaller, error) {
	contract, err := bindVRFCoordinatorV2EventsMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockCaller{contract: contract}, nil
}

func NewVRFCoordinatorV2EventsMockTransactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorV2EventsMockTransactor, error) {
	contract, err := bindVRFCoordinatorV2EventsMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockTransactor{contract: contract}, nil
}

func NewVRFCoordinatorV2EventsMockFilterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorV2EventsMockFilterer, error) {
	contract, err := bindVRFCoordinatorV2EventsMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockFilterer{contract: contract}, nil
}

func bindVRFCoordinatorV2EventsMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorV2EventsMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2EventsMock.Contract.VRFCoordinatorV2EventsMockCaller.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.VRFCoordinatorV2EventsMockTransactor.contract.Transfer(opts)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.VRFCoordinatorV2EventsMockTransactor.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2EventsMock.Contract.contract.Call(opts, result, method, params...)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.contract.Transfer(opts)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.contract.Transact(opts, method, params...)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitConfigSet(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2EventsMockFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitConfigSet", minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitConfigSet(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2EventsMockFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitConfigSet(&_VRFCoordinatorV2EventsMock.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitConfigSet(minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2EventsMockFeeConfig) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitConfigSet(&_VRFCoordinatorV2EventsMock.TransactOpts, minimumRequestConfirmations, maxGasLimit, stalenessSeconds, gasAfterPaymentCalculation, fallbackWeiPerUnitLink, feeConfig)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitFundsRecovered(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitFundsRecovered", to, amount)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitFundsRecovered(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitFundsRecovered(&_VRFCoordinatorV2EventsMock.TransactOpts, to, amount)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitFundsRecovered(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitFundsRecovered(&_VRFCoordinatorV2EventsMock.TransactOpts, to, amount)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitOwnershipTransferRequested(&_VRFCoordinatorV2EventsMock.TransactOpts, from, to)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitOwnershipTransferRequested(&_VRFCoordinatorV2EventsMock.TransactOpts, from, to)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitOwnershipTransferred(&_VRFCoordinatorV2EventsMock.TransactOpts, from, to)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitOwnershipTransferred(&_VRFCoordinatorV2EventsMock.TransactOpts, from, to)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitProvingKeyDeregistered(opts *bind.TransactOpts, keyHash [32]byte, oracle common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitProvingKeyDeregistered", keyHash, oracle)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitProvingKeyDeregistered(keyHash [32]byte, oracle common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitProvingKeyDeregistered(&_VRFCoordinatorV2EventsMock.TransactOpts, keyHash, oracle)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitProvingKeyDeregistered(keyHash [32]byte, oracle common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitProvingKeyDeregistered(&_VRFCoordinatorV2EventsMock.TransactOpts, keyHash, oracle)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitProvingKeyRegistered(opts *bind.TransactOpts, keyHash [32]byte, oracle common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitProvingKeyRegistered", keyHash, oracle)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitProvingKeyRegistered(keyHash [32]byte, oracle common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitProvingKeyRegistered(&_VRFCoordinatorV2EventsMock.TransactOpts, keyHash, oracle)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitProvingKeyRegistered(keyHash [32]byte, oracle common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitProvingKeyRegistered(&_VRFCoordinatorV2EventsMock.TransactOpts, keyHash, oracle)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitRandomWordsFulfilled(opts *bind.TransactOpts, requestId *big.Int, outputSeed *big.Int, payment *big.Int, success bool) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitRandomWordsFulfilled", requestId, outputSeed, payment, success)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitRandomWordsFulfilled(requestId *big.Int, outputSeed *big.Int, payment *big.Int, success bool) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitRandomWordsFulfilled(&_VRFCoordinatorV2EventsMock.TransactOpts, requestId, outputSeed, payment, success)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitRandomWordsFulfilled(requestId *big.Int, outputSeed *big.Int, payment *big.Int, success bool) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitRandomWordsFulfilled(&_VRFCoordinatorV2EventsMock.TransactOpts, requestId, outputSeed, payment, success)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitRandomWordsRequested(opts *bind.TransactOpts, keyHash [32]byte, requestId *big.Int, preSeed *big.Int, subId uint64, minimumRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32, sender common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitRandomWordsRequested", keyHash, requestId, preSeed, subId, minimumRequestConfirmations, callbackGasLimit, numWords, sender)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitRandomWordsRequested(keyHash [32]byte, requestId *big.Int, preSeed *big.Int, subId uint64, minimumRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32, sender common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitRandomWordsRequested(&_VRFCoordinatorV2EventsMock.TransactOpts, keyHash, requestId, preSeed, subId, minimumRequestConfirmations, callbackGasLimit, numWords, sender)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitRandomWordsRequested(keyHash [32]byte, requestId *big.Int, preSeed *big.Int, subId uint64, minimumRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32, sender common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitRandomWordsRequested(&_VRFCoordinatorV2EventsMock.TransactOpts, keyHash, requestId, preSeed, subId, minimumRequestConfirmations, callbackGasLimit, numWords, sender)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitSubscriptionCanceled(opts *bind.TransactOpts, subId uint64, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitSubscriptionCanceled", subId, to, amount)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitSubscriptionCanceled(subId uint64, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionCanceled(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, to, amount)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitSubscriptionCanceled(subId uint64, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionCanceled(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, to, amount)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitSubscriptionConsumerAdded(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitSubscriptionConsumerAdded", subId, consumer)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitSubscriptionConsumerAdded(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionConsumerAdded(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitSubscriptionConsumerAdded(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionConsumerAdded(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitSubscriptionConsumerRemoved(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitSubscriptionConsumerRemoved", subId, consumer)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitSubscriptionConsumerRemoved(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionConsumerRemoved(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitSubscriptionConsumerRemoved(subId uint64, consumer common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionConsumerRemoved(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, consumer)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitSubscriptionCreated(opts *bind.TransactOpts, subId uint64, owner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitSubscriptionCreated", subId, owner)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitSubscriptionCreated(subId uint64, owner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionCreated(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, owner)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitSubscriptionCreated(subId uint64, owner common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionCreated(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, owner)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitSubscriptionFunded(opts *bind.TransactOpts, subId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitSubscriptionFunded", subId, oldBalance, newBalance)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitSubscriptionFunded(subId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionFunded(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, oldBalance, newBalance)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitSubscriptionFunded(subId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionFunded(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, oldBalance, newBalance)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitSubscriptionOwnerTransferRequested(opts *bind.TransactOpts, subId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitSubscriptionOwnerTransferRequested", subId, from, to)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitSubscriptionOwnerTransferRequested(subId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionOwnerTransferRequested(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, from, to)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitSubscriptionOwnerTransferRequested(subId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionOwnerTransferRequested(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, from, to)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactor) EmitSubscriptionOwnerTransferred(opts *bind.TransactOpts, subId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.contract.Transact(opts, "emitSubscriptionOwnerTransferred", subId, from, to)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockSession) EmitSubscriptionOwnerTransferred(subId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionOwnerTransferred(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, from, to)
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockTransactorSession) EmitSubscriptionOwnerTransferred(subId uint64, from common.Address, to common.Address) (*types.Transaction, error) {
	return _VRFCoordinatorV2EventsMock.Contract.EmitSubscriptionOwnerTransferred(&_VRFCoordinatorV2EventsMock.TransactOpts, subId, from, to)
}

type VRFCoordinatorV2EventsMockConfigSetIterator struct {
	Event *VRFCoordinatorV2EventsMockConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockConfigSet)
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
		it.Event = new(VRFCoordinatorV2EventsMockConfigSet)
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

func (it *VRFCoordinatorV2EventsMockConfigSetIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockConfigSet struct {
	MinimumRequestConfirmations uint16
	MaxGasLimit                 uint32
	StalenessSeconds            uint32
	GasAfterPaymentCalculation  uint32
	FallbackWeiPerUnitLink      *big.Int
	FeeConfig                   VRFCoordinatorV2EventsMockFeeConfig
	Raw                         types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV2EventsMockConfigSetIterator, error) {

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockConfigSetIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockConfigSet) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockConfigSet)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseConfigSet(log types.Log) (*VRFCoordinatorV2EventsMockConfigSet, error) {
	event := new(VRFCoordinatorV2EventsMockConfigSet)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockFundsRecoveredIterator struct {
	Event *VRFCoordinatorV2EventsMockFundsRecovered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockFundsRecoveredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockFundsRecovered)
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
		it.Event = new(VRFCoordinatorV2EventsMockFundsRecovered)
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

func (it *VRFCoordinatorV2EventsMockFundsRecoveredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockFundsRecoveredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockFundsRecovered struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2EventsMockFundsRecoveredIterator, error) {

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockFundsRecoveredIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "FundsRecovered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockFundsRecovered) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "FundsRecovered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockFundsRecovered)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseFundsRecovered(log types.Log) (*VRFCoordinatorV2EventsMockFundsRecovered, error) {
	event := new(VRFCoordinatorV2EventsMockFundsRecovered)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "FundsRecovered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockOwnershipTransferRequestedIterator struct {
	Event *VRFCoordinatorV2EventsMockOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockOwnershipTransferRequested)
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
		it.Event = new(VRFCoordinatorV2EventsMockOwnershipTransferRequested)
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

func (it *VRFCoordinatorV2EventsMockOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2EventsMockOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockOwnershipTransferRequestedIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockOwnershipTransferRequested)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV2EventsMockOwnershipTransferRequested, error) {
	event := new(VRFCoordinatorV2EventsMockOwnershipTransferRequested)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockOwnershipTransferredIterator struct {
	Event *VRFCoordinatorV2EventsMockOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockOwnershipTransferred)
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
		it.Event = new(VRFCoordinatorV2EventsMockOwnershipTransferred)
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

func (it *VRFCoordinatorV2EventsMockOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2EventsMockOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockOwnershipTransferredIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockOwnershipTransferred)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV2EventsMockOwnershipTransferred, error) {
	event := new(VRFCoordinatorV2EventsMockOwnershipTransferred)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockProvingKeyDeregisteredIterator struct {
	Event *VRFCoordinatorV2EventsMockProvingKeyDeregistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockProvingKeyDeregisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockProvingKeyDeregistered)
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
		it.Event = new(VRFCoordinatorV2EventsMockProvingKeyDeregistered)
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

func (it *VRFCoordinatorV2EventsMockProvingKeyDeregisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockProvingKeyDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockProvingKeyDeregistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2EventsMockProvingKeyDeregisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockProvingKeyDeregisteredIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "ProvingKeyDeregistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "ProvingKeyDeregistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockProvingKeyDeregistered)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV2EventsMockProvingKeyDeregistered, error) {
	event := new(VRFCoordinatorV2EventsMockProvingKeyDeregistered)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "ProvingKeyDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockProvingKeyRegisteredIterator struct {
	Event *VRFCoordinatorV2EventsMockProvingKeyRegistered

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockProvingKeyRegisteredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockProvingKeyRegistered)
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
		it.Event = new(VRFCoordinatorV2EventsMockProvingKeyRegistered)
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

func (it *VRFCoordinatorV2EventsMockProvingKeyRegisteredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockProvingKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockProvingKeyRegistered struct {
	KeyHash [32]byte
	Oracle  common.Address
	Raw     types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2EventsMockProvingKeyRegisteredIterator, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockProvingKeyRegisteredIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "ProvingKeyRegistered", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockProvingKeyRegistered, oracle []common.Address) (event.Subscription, error) {

	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "ProvingKeyRegistered", oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockProvingKeyRegistered)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV2EventsMockProvingKeyRegistered, error) {
	event := new(VRFCoordinatorV2EventsMockProvingKeyRegistered)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "ProvingKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockRandomWordsFulfilledIterator struct {
	Event *VRFCoordinatorV2EventsMockRandomWordsFulfilled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockRandomWordsFulfilledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockRandomWordsFulfilled)
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
		it.Event = new(VRFCoordinatorV2EventsMockRandomWordsFulfilled)
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

func (it *VRFCoordinatorV2EventsMockRandomWordsFulfilledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockRandomWordsFulfilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockRandomWordsFulfilled struct {
	RequestId  *big.Int
	OutputSeed *big.Int
	Payment    *big.Int
	Success    bool
	Raw        types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFCoordinatorV2EventsMockRandomWordsFulfilledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockRandomWordsFulfilledIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "RandomWordsFulfilled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockRandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "RandomWordsFulfilled", requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockRandomWordsFulfilled)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV2EventsMockRandomWordsFulfilled, error) {
	event := new(VRFCoordinatorV2EventsMockRandomWordsFulfilled)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "RandomWordsFulfilled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockRandomWordsRequestedIterator struct {
	Event *VRFCoordinatorV2EventsMockRandomWordsRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockRandomWordsRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockRandomWordsRequested)
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
		it.Event = new(VRFCoordinatorV2EventsMockRandomWordsRequested)
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

func (it *VRFCoordinatorV2EventsMockRandomWordsRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockRandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockRandomWordsRequested struct {
	KeyHash                     [32]byte
	RequestId                   *big.Int
	PreSeed                     *big.Int
	SubId                       uint64
	MinimumRequestConfirmations uint16
	CallbackGasLimit            uint32
	NumWords                    uint32
	Sender                      common.Address
	Raw                         types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*VRFCoordinatorV2EventsMockRandomWordsRequestedIterator, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockRandomWordsRequestedIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockRandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error) {

	var keyHashRule []interface{}
	for _, keyHashItem := range keyHash {
		keyHashRule = append(keyHashRule, keyHashItem)
	}

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "RandomWordsRequested", keyHashRule, subIdRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockRandomWordsRequested)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV2EventsMockRandomWordsRequested, error) {
	event := new(VRFCoordinatorV2EventsMockRandomWordsRequested)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockSubscriptionCanceledIterator struct {
	Event *VRFCoordinatorV2EventsMockSubscriptionCanceled

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockSubscriptionCanceledIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockSubscriptionCanceled)
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
		it.Event = new(VRFCoordinatorV2EventsMockSubscriptionCanceled)
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

func (it *VRFCoordinatorV2EventsMockSubscriptionCanceledIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockSubscriptionCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockSubscriptionCanceled struct {
	SubId  uint64
	To     common.Address
	Amount *big.Int
	Raw    types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionCanceledIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockSubscriptionCanceledIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "SubscriptionCanceled", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionCanceled, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "SubscriptionCanceled", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockSubscriptionCanceled)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionCanceled, error) {
	event := new(VRFCoordinatorV2EventsMockSubscriptionCanceled)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockSubscriptionConsumerAddedIterator struct {
	Event *VRFCoordinatorV2EventsMockSubscriptionConsumerAdded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockSubscriptionConsumerAddedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockSubscriptionConsumerAdded)
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
		it.Event = new(VRFCoordinatorV2EventsMockSubscriptionConsumerAdded)
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

func (it *VRFCoordinatorV2EventsMockSubscriptionConsumerAddedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockSubscriptionConsumerAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockSubscriptionConsumerAdded struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionConsumerAddedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockSubscriptionConsumerAddedIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "SubscriptionConsumerAdded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionConsumerAdded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "SubscriptionConsumerAdded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockSubscriptionConsumerAdded)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionConsumerAdded, error) {
	event := new(VRFCoordinatorV2EventsMockSubscriptionConsumerAdded)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionConsumerAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockSubscriptionConsumerRemovedIterator struct {
	Event *VRFCoordinatorV2EventsMockSubscriptionConsumerRemoved

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockSubscriptionConsumerRemovedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockSubscriptionConsumerRemoved)
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
		it.Event = new(VRFCoordinatorV2EventsMockSubscriptionConsumerRemoved)
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

func (it *VRFCoordinatorV2EventsMockSubscriptionConsumerRemovedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockSubscriptionConsumerRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockSubscriptionConsumerRemoved struct {
	SubId    uint64
	Consumer common.Address
	Raw      types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionConsumerRemovedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockSubscriptionConsumerRemovedIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "SubscriptionConsumerRemoved", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "SubscriptionConsumerRemoved", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockSubscriptionConsumerRemoved)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionConsumerRemoved, error) {
	event := new(VRFCoordinatorV2EventsMockSubscriptionConsumerRemoved)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionConsumerRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockSubscriptionCreatedIterator struct {
	Event *VRFCoordinatorV2EventsMockSubscriptionCreated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockSubscriptionCreatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockSubscriptionCreated)
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
		it.Event = new(VRFCoordinatorV2EventsMockSubscriptionCreated)
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

func (it *VRFCoordinatorV2EventsMockSubscriptionCreatedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockSubscriptionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockSubscriptionCreated struct {
	SubId uint64
	Owner common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionCreatedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockSubscriptionCreatedIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "SubscriptionCreated", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionCreated, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "SubscriptionCreated", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockSubscriptionCreated)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionCreated, error) {
	event := new(VRFCoordinatorV2EventsMockSubscriptionCreated)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockSubscriptionFundedIterator struct {
	Event *VRFCoordinatorV2EventsMockSubscriptionFunded

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockSubscriptionFundedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockSubscriptionFunded)
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
		it.Event = new(VRFCoordinatorV2EventsMockSubscriptionFunded)
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

func (it *VRFCoordinatorV2EventsMockSubscriptionFundedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockSubscriptionFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockSubscriptionFunded struct {
	SubId      uint64
	OldBalance *big.Int
	NewBalance *big.Int
	Raw        types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionFundedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockSubscriptionFundedIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "SubscriptionFunded", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionFunded, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "SubscriptionFunded", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockSubscriptionFunded)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionFunded, error) {
	event := new(VRFCoordinatorV2EventsMockSubscriptionFunded)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequestedIterator struct {
	Event *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequested)
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
		it.Event = new(VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequested)
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

func (it *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequested struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequestedIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequestedIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "SubscriptionOwnerTransferRequested", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "SubscriptionOwnerTransferRequested", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequested)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequested, error) {
	event := new(VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequested)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type VRFCoordinatorV2EventsMockSubscriptionOwnerTransferredIterator struct {
	Event *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2EventsMockSubscriptionOwnerTransferred)
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
		it.Event = new(VRFCoordinatorV2EventsMockSubscriptionOwnerTransferred)
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

func (it *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferredIterator) Error() error {
	return it.fail
}

func (it *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type VRFCoordinatorV2EventsMockSubscriptionOwnerTransferred struct {
	SubId uint64
	From  common.Address
	To    common.Address
	Raw   types.Log
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionOwnerTransferredIterator, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.FilterLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2EventsMockSubscriptionOwnerTransferredIterator{contract: _VRFCoordinatorV2EventsMock.contract, event: "SubscriptionOwnerTransferred", logs: logs, sub: sub}, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error) {

	var subIdRule []interface{}
	for _, subIdItem := range subId {
		subIdRule = append(subIdRule, subIdItem)
	}

	logs, sub, err := _VRFCoordinatorV2EventsMock.contract.WatchLogs(opts, "SubscriptionOwnerTransferred", subIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(VRFCoordinatorV2EventsMockSubscriptionOwnerTransferred)
				if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
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

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMockFilterer) ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionOwnerTransferred, error) {
	event := new(VRFCoordinatorV2EventsMockSubscriptionOwnerTransferred)
	if err := _VRFCoordinatorV2EventsMock.contract.UnpackLog(event, "SubscriptionOwnerTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _VRFCoordinatorV2EventsMock.abi.Events["ConfigSet"].ID:
		return _VRFCoordinatorV2EventsMock.ParseConfigSet(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["FundsRecovered"].ID:
		return _VRFCoordinatorV2EventsMock.ParseFundsRecovered(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["OwnershipTransferRequested"].ID:
		return _VRFCoordinatorV2EventsMock.ParseOwnershipTransferRequested(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["OwnershipTransferred"].ID:
		return _VRFCoordinatorV2EventsMock.ParseOwnershipTransferred(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["ProvingKeyDeregistered"].ID:
		return _VRFCoordinatorV2EventsMock.ParseProvingKeyDeregistered(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["ProvingKeyRegistered"].ID:
		return _VRFCoordinatorV2EventsMock.ParseProvingKeyRegistered(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["RandomWordsFulfilled"].ID:
		return _VRFCoordinatorV2EventsMock.ParseRandomWordsFulfilled(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["RandomWordsRequested"].ID:
		return _VRFCoordinatorV2EventsMock.ParseRandomWordsRequested(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["SubscriptionCanceled"].ID:
		return _VRFCoordinatorV2EventsMock.ParseSubscriptionCanceled(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["SubscriptionConsumerAdded"].ID:
		return _VRFCoordinatorV2EventsMock.ParseSubscriptionConsumerAdded(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["SubscriptionConsumerRemoved"].ID:
		return _VRFCoordinatorV2EventsMock.ParseSubscriptionConsumerRemoved(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["SubscriptionCreated"].ID:
		return _VRFCoordinatorV2EventsMock.ParseSubscriptionCreated(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["SubscriptionFunded"].ID:
		return _VRFCoordinatorV2EventsMock.ParseSubscriptionFunded(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["SubscriptionOwnerTransferRequested"].ID:
		return _VRFCoordinatorV2EventsMock.ParseSubscriptionOwnerTransferRequested(log)
	case _VRFCoordinatorV2EventsMock.abi.Events["SubscriptionOwnerTransferred"].ID:
		return _VRFCoordinatorV2EventsMock.ParseSubscriptionOwnerTransferred(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (VRFCoordinatorV2EventsMockConfigSet) Topic() common.Hash {
	return common.HexToHash("0xc21e3bd2e0b339d2848f0dd956947a88966c242c0c0c582a33137a5c1ceb5cb2")
}

func (VRFCoordinatorV2EventsMockFundsRecovered) Topic() common.Hash {
	return common.HexToHash("0x59bfc682b673f8cbf945f1e454df9334834abf7dfe7f92237ca29ecb9b436600")
}

func (VRFCoordinatorV2EventsMockOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (VRFCoordinatorV2EventsMockOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (VRFCoordinatorV2EventsMockProvingKeyDeregistered) Topic() common.Hash {
	return common.HexToHash("0x72be339577868f868798bac2c93e52d6f034fef4689a9848996c14ebb7416c0d")
}

func (VRFCoordinatorV2EventsMockProvingKeyRegistered) Topic() common.Hash {
	return common.HexToHash("0xe729ae16526293f74ade739043022254f1489f616295a25bf72dfb4511ed73b8")
}

func (VRFCoordinatorV2EventsMockRandomWordsFulfilled) Topic() common.Hash {
	return common.HexToHash("0x7dffc5ae5ee4e2e4df1651cf6ad329a73cebdb728f37ea0187b9b17e036756e4")
}

func (VRFCoordinatorV2EventsMockRandomWordsRequested) Topic() common.Hash {
	return common.HexToHash("0x63373d1c4696214b898952999c9aaec57dac1ee2723cec59bea6888f489a9772")
}

func (VRFCoordinatorV2EventsMockSubscriptionCanceled) Topic() common.Hash {
	return common.HexToHash("0xe8ed5b475a5b5987aa9165e8731bb78043f39eee32ec5a1169a89e27fcd49815")
}

func (VRFCoordinatorV2EventsMockSubscriptionConsumerAdded) Topic() common.Hash {
	return common.HexToHash("0x43dc749a04ac8fb825cbd514f7c0e13f13bc6f2ee66043b76629d51776cff8e0")
}

func (VRFCoordinatorV2EventsMockSubscriptionConsumerRemoved) Topic() common.Hash {
	return common.HexToHash("0x182bff9831466789164ca77075fffd84916d35a8180ba73c27e45634549b445b")
}

func (VRFCoordinatorV2EventsMockSubscriptionCreated) Topic() common.Hash {
	return common.HexToHash("0x464722b4166576d3dcbba877b999bc35cf911f4eaf434b7eba68fa113951d0bf")
}

func (VRFCoordinatorV2EventsMockSubscriptionFunded) Topic() common.Hash {
	return common.HexToHash("0xd39ec07f4e209f627a4c427971473820dc129761ba28de8906bd56f57101d4f8")
}

func (VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x69436ea6df009049404f564eff6622cd00522b0bd6a89efd9e52a355c4a879be")
}

func (VRFCoordinatorV2EventsMockSubscriptionOwnerTransferred) Topic() common.Hash {
	return common.HexToHash("0x6f1dc65165ffffedfd8e507b4a0f1fcfdada045ed11f6c26ba27cedfe87802f0")
}

func (_VRFCoordinatorV2EventsMock *VRFCoordinatorV2EventsMock) Address() common.Address {
	return _VRFCoordinatorV2EventsMock.address
}

type VRFCoordinatorV2EventsMockInterface interface {
	EmitConfigSet(opts *bind.TransactOpts, minimumRequestConfirmations uint16, maxGasLimit uint32, stalenessSeconds uint32, gasAfterPaymentCalculation uint32, fallbackWeiPerUnitLink *big.Int, feeConfig VRFCoordinatorV2EventsMockFeeConfig) (*types.Transaction, error)

	EmitFundsRecovered(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error)

	EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitProvingKeyDeregistered(opts *bind.TransactOpts, keyHash [32]byte, oracle common.Address) (*types.Transaction, error)

	EmitProvingKeyRegistered(opts *bind.TransactOpts, keyHash [32]byte, oracle common.Address) (*types.Transaction, error)

	EmitRandomWordsFulfilled(opts *bind.TransactOpts, requestId *big.Int, outputSeed *big.Int, payment *big.Int, success bool) (*types.Transaction, error)

	EmitRandomWordsRequested(opts *bind.TransactOpts, keyHash [32]byte, requestId *big.Int, preSeed *big.Int, subId uint64, minimumRequestConfirmations uint16, callbackGasLimit uint32, numWords uint32, sender common.Address) (*types.Transaction, error)

	EmitSubscriptionCanceled(opts *bind.TransactOpts, subId uint64, to common.Address, amount *big.Int) (*types.Transaction, error)

	EmitSubscriptionConsumerAdded(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)

	EmitSubscriptionConsumerRemoved(opts *bind.TransactOpts, subId uint64, consumer common.Address) (*types.Transaction, error)

	EmitSubscriptionCreated(opts *bind.TransactOpts, subId uint64, owner common.Address) (*types.Transaction, error)

	EmitSubscriptionFunded(opts *bind.TransactOpts, subId uint64, oldBalance *big.Int, newBalance *big.Int) (*types.Transaction, error)

	EmitSubscriptionOwnerTransferRequested(opts *bind.TransactOpts, subId uint64, from common.Address, to common.Address) (*types.Transaction, error)

	EmitSubscriptionOwnerTransferred(opts *bind.TransactOpts, subId uint64, from common.Address, to common.Address) (*types.Transaction, error)

	FilterConfigSet(opts *bind.FilterOpts) (*VRFCoordinatorV2EventsMockConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*VRFCoordinatorV2EventsMockConfigSet, error)

	FilterFundsRecovered(opts *bind.FilterOpts) (*VRFCoordinatorV2EventsMockFundsRecoveredIterator, error)

	WatchFundsRecovered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockFundsRecovered) (event.Subscription, error)

	ParseFundsRecovered(log types.Log) (*VRFCoordinatorV2EventsMockFundsRecovered, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2EventsMockOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*VRFCoordinatorV2EventsMockOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*VRFCoordinatorV2EventsMockOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*VRFCoordinatorV2EventsMockOwnershipTransferred, error)

	FilterProvingKeyDeregistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2EventsMockProvingKeyDeregisteredIterator, error)

	WatchProvingKeyDeregistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockProvingKeyDeregistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyDeregistered(log types.Log) (*VRFCoordinatorV2EventsMockProvingKeyDeregistered, error)

	FilterProvingKeyRegistered(opts *bind.FilterOpts, oracle []common.Address) (*VRFCoordinatorV2EventsMockProvingKeyRegisteredIterator, error)

	WatchProvingKeyRegistered(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockProvingKeyRegistered, oracle []common.Address) (event.Subscription, error)

	ParseProvingKeyRegistered(log types.Log) (*VRFCoordinatorV2EventsMockProvingKeyRegistered, error)

	FilterRandomWordsFulfilled(opts *bind.FilterOpts, requestId []*big.Int) (*VRFCoordinatorV2EventsMockRandomWordsFulfilledIterator, error)

	WatchRandomWordsFulfilled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockRandomWordsFulfilled, requestId []*big.Int) (event.Subscription, error)

	ParseRandomWordsFulfilled(log types.Log) (*VRFCoordinatorV2EventsMockRandomWordsFulfilled, error)

	FilterRandomWordsRequested(opts *bind.FilterOpts, keyHash [][32]byte, subId []uint64, sender []common.Address) (*VRFCoordinatorV2EventsMockRandomWordsRequestedIterator, error)

	WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockRandomWordsRequested, keyHash [][32]byte, subId []uint64, sender []common.Address) (event.Subscription, error)

	ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV2EventsMockRandomWordsRequested, error)

	FilterSubscriptionCanceled(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionCanceledIterator, error)

	WatchSubscriptionCanceled(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionCanceled, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCanceled(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionCanceled, error)

	FilterSubscriptionConsumerAdded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionConsumerAddedIterator, error)

	WatchSubscriptionConsumerAdded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionConsumerAdded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerAdded(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionConsumerAdded, error)

	FilterSubscriptionConsumerRemoved(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionConsumerRemovedIterator, error)

	WatchSubscriptionConsumerRemoved(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionConsumerRemoved, subId []uint64) (event.Subscription, error)

	ParseSubscriptionConsumerRemoved(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionConsumerRemoved, error)

	FilterSubscriptionCreated(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionCreatedIterator, error)

	WatchSubscriptionCreated(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionCreated, subId []uint64) (event.Subscription, error)

	ParseSubscriptionCreated(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionCreated, error)

	FilterSubscriptionFunded(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionFundedIterator, error)

	WatchSubscriptionFunded(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionFunded, subId []uint64) (event.Subscription, error)

	ParseSubscriptionFunded(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionFunded, error)

	FilterSubscriptionOwnerTransferRequested(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequestedIterator, error)

	WatchSubscriptionOwnerTransferRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequested, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferRequested(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionOwnerTransferRequested, error)

	FilterSubscriptionOwnerTransferred(opts *bind.FilterOpts, subId []uint64) (*VRFCoordinatorV2EventsMockSubscriptionOwnerTransferredIterator, error)

	WatchSubscriptionOwnerTransferred(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2EventsMockSubscriptionOwnerTransferred, subId []uint64) (event.Subscription, error)

	ParseSubscriptionOwnerTransferred(log types.Log) (*VRFCoordinatorV2EventsMockSubscriptionOwnerTransferred, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
