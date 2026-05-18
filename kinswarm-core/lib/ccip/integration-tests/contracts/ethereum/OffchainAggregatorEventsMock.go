// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethereum

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

var OffchainAggregatorEventsMockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"BillingAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"BillingSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"ConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_oldLinkToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_newLinkToken\",\"type\":\"address\"}],\"name\":\"LinkTokenSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"NewTransmission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"OraclePaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"PayeeshipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"RequesterAccessControllerSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"}],\"name\":\"RoundRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"previousGasLimit\",\"type\":\"uint32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"currentValidator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"currentGasLimit\",\"type\":\"uint32\"}],\"name\":\"ValidatorConfigSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"emitAnswerUpdated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"emitBillingAccessControllerSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maximumGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"reasonableGasPrice\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"microLinkPerEth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerObservation\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"linkGweiPerTransmission\",\"type\":\"uint32\"}],\"name\":\"emitBillingSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"previousConfigBlockNumber\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"configCount\",\"type\":\"uint64\"},{\"internalType\":\"address[]\",\"name\":\"signers\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"transmitters\",\"type\":\"address[]\"},{\"internalType\":\"uint8\",\"name\":\"threshold\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"encodedConfigVersion\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"emitConfigSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oldLinkToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_newLinkToken\",\"type\":\"address\"}],\"name\":\"emitLinkTokenSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"emitNewRound\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"aggregatorRoundId\",\"type\":\"uint32\"},{\"internalType\":\"int192\",\"name\":\"answer\",\"type\":\"int192\"},{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"int192[]\",\"name\":\"observations\",\"type\":\"int192[]\"},{\"internalType\":\"bytes\",\"name\":\"observers\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"rawReportContext\",\"type\":\"bytes32\"}],\"name\":\"emitNewTransmission\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"payee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"linkToken\",\"type\":\"address\"}],\"name\":\"emitOraclePaid\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"emitOwnershipTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"emitPayeeshipTransferRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transmitter\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"emitPayeeshipTransferred\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"old\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"}],\"name\":\"emitRequesterAccessControllerSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"internalType\":\"bytes16\",\"name\":\"configDigest\",\"type\":\"bytes16\"},{\"internalType\":\"uint32\",\"name\":\"epoch\",\"type\":\"uint32\"},{\"internalType\":\"uint8\",\"name\":\"round\",\"type\":\"uint8\"}],\"name\":\"emitRoundRequested\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"previousValidator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"previousGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"currentValidator\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"currentGasLimit\",\"type\":\"uint32\"}],\"name\":\"emitValidatorConfigSet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610faf806100206000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c80638b28369d11610097578063a57c1e0c11610066578063a57c1e0c146101cd578063b019b4e8146101e0578063f7420bc2146101f3578063faf1347c1461020657600080fd5b80638b28369d1461018157806395aa14461461019457806399e1a39b146101a7578063a3296557146101ba57600080fd5b806358e1e734116100d357806358e1e734146101355780636602e6ce14610148578063715bd44e1461015b57806389ffde8d1461016e57600080fd5b8063275c7ea4146100fa5780632c769fd71461010f578063448be1e014610122575b600080fd5b61010d610108366004610c76565b610219565b005b61010d61011d366004610aad565b610265565b61010d610130366004610931565b6102a5565b61010d610143366004610964565b6102fa565b61010d610156366004610a64565b610370565b61010d6101693660046109f4565b6103d3565b61010d61017c366004610931565b61045b565b61010d61018f366004610ad9565b6104a8565b61010d6101a2366004610b0e565b6104f1565b61010d6101b5366004610964565b61053f565b61010d6101c8366004610c11565b6105b5565b61010d6101db366004610931565b610614565b61010d6101ee366004610931565b610672565b61010d610201366004610931565b6106d0565b61010d6102143660046109a7565b61072e565b7f25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb9878787878787876040516102549796959493929190610e8a565b60405180910390a150505050505050565b81837f0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f8360405161029891815260200190565b60405180910390a3505050565b6040805173ffffffffffffffffffffffffffffffffffffffff8085168252831660208201527f793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d4891291015b60405180910390a15050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b360405160405180910390a4505050565b6040805163ffffffff80861682528316602082015273ffffffffffffffffffffffffffffffffffffffff80851692908716917fb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541910160405180910390a350505050565b604080517fffffffffffffffffffffffffffffffff000000000000000000000000000000008516815263ffffffff8416602082015260ff831681830152905173ffffffffffffffffffffffffffffffffffffffff8616917f3ea16a923ff4b1df6526e854c9e3a995c43385d70e73359e10623c74f0b52037919081900360600190a250505050565b6040805173ffffffffffffffffffffffffffffffffffffffff8085168252831660208201527f27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae63491016102ee565b8173ffffffffffffffffffffffffffffffffffffffff16837f0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac602718360405161029891815260200190565b8563ffffffff167ff6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451868686868660405161052f959493929190610dfd565b60405180910390a2505050505050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e3836760405160405180910390a4505050565b6040805163ffffffff878116825286811660208301528581168284015284811660608301528316608082015290517fd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b69181900360a00190a15050505050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a60405160405180910390a35050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae127860405160405180910390a35050565b8073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167fd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c856040516107a491815260200190565b60405180910390a450505050565b803573ffffffffffffffffffffffffffffffffffffffff811681146107d657600080fd5b919050565b600082601f8301126107ec57600080fd5b813560206108016107fc83610f4f565b610f00565b80838252828201915082860187848660051b890101111561082157600080fd5b60005b8581101561084757610835826107b2565b84529284019290840190600101610824565b5090979650505050505050565b600082601f83011261086557600080fd5b813567ffffffffffffffff81111561087f5761087f610f73565b6108b060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f84011601610f00565b8181528460208386010111156108c557600080fd5b816020850160208301376000918101602001919091529392505050565b8035601781900b81146107d657600080fd5b803563ffffffff811681146107d657600080fd5b803567ffffffffffffffff811681146107d657600080fd5b803560ff811681146107d657600080fd5b6000806040838503121561094457600080fd5b61094d836107b2565b915061095b602084016107b2565b90509250929050565b60008060006060848603121561097957600080fd5b610982846107b2565b9250610990602085016107b2565b915061099e604085016107b2565b90509250925092565b600080600080608085870312156109bd57600080fd5b6109c6856107b2565b93506109d4602086016107b2565b9250604085013591506109e9606086016107b2565b905092959194509250565b60008060008060808587031215610a0a57600080fd5b610a13856107b2565b935060208501357fffffffffffffffffffffffffffffffff0000000000000000000000000000000081168114610a4857600080fd5b9250610a56604086016108f4565b91506109e960608601610920565b60008060008060808587031215610a7a57600080fd5b610a83856107b2565b9350610a91602086016108f4565b9250610a9f604086016107b2565b91506109e9606086016108f4565b600080600060608486031215610ac257600080fd5b505081359360208301359350604090920135919050565b600080600060608486031215610aee57600080fd5b83359250610afe602085016107b2565b9150604084013590509250925092565b60008060008060008060c08789031215610b2757600080fd5b610b30876108f4565b95506020610b3f8189016108e2565b9550610b4d604089016107b2565b9450606088013567ffffffffffffffff80821115610b6a57600080fd5b818a0191508a601f830112610b7e57600080fd5b8135610b8c6107fc82610f4f565b8082825285820191508585018e878560051b8801011115610bac57600080fd5b600095505b83861015610bd657610bc2816108e2565b835260019590950194918601918601610bb1565b509750505060808a0135925080831115610bef57600080fd5b5050610bfd89828a01610854565b92505060a087013590509295509295509295565b600080600080600060a08688031215610c2957600080fd5b610c32866108f4565b9450610c40602087016108f4565b9350610c4e604087016108f4565b9250610c5c606087016108f4565b9150610c6a608087016108f4565b90509295509295909350565b600080600080600080600060e0888a031215610c9157600080fd5b610c9a886108f4565b9650610ca860208901610908565b9550604088013567ffffffffffffffff80821115610cc557600080fd5b610cd18b838c016107db565b965060608a0135915080821115610ce757600080fd5b610cf38b838c016107db565b9550610d0160808b01610920565b9450610d0f60a08b01610908565b935060c08a0135915080821115610d2557600080fd5b50610d328a828b01610854565b91505092959891949750929550565b600081518084526020808501945080840160005b83811015610d8757815173ffffffffffffffffffffffffffffffffffffffff1687529582019590820190600101610d55565b509495945050505050565b6000815180845260005b81811015610db857602081850181015186830182015201610d9c565b81811115610dca576000602083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b600060a08201601788810b8452602073ffffffffffffffffffffffffffffffffffffffff89168186015260a0604086015282885180855260c087019150828a01945060005b81811015610e60578551850b83529483019491830191600101610e42565b50508581036060870152610e748189610d92565b9450505050508260808301529695505050505050565b63ffffffff88168152600067ffffffffffffffff808916602084015260e06040840152610eba60e0840189610d41565b8381036060850152610ecc8189610d41565b905060ff8716608085015281861660a085015283810360c0850152610ef18186610d92565b9b9a5050505050505050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715610f4757610f47610f73565b604052919050565b600067ffffffffffffffff821115610f6957610f69610f73565b5060051b60200190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a",
}

var OffchainAggregatorEventsMockABI = OffchainAggregatorEventsMockMetaData.ABI

var OffchainAggregatorEventsMockBin = OffchainAggregatorEventsMockMetaData.Bin

func DeployOffchainAggregatorEventsMock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OffchainAggregatorEventsMock, error) {
	parsed, err := OffchainAggregatorEventsMockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OffchainAggregatorEventsMockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OffchainAggregatorEventsMock{OffchainAggregatorEventsMockCaller: OffchainAggregatorEventsMockCaller{contract: contract}, OffchainAggregatorEventsMockTransactor: OffchainAggregatorEventsMockTransactor{contract: contract}, OffchainAggregatorEventsMockFilterer: OffchainAggregatorEventsMockFilterer{contract: contract}}, nil
}

type OffchainAggregatorEventsMock struct {
	address common.Address
	abi     abi.ABI
	OffchainAggregatorEventsMockCaller
	OffchainAggregatorEventsMockTransactor
	OffchainAggregatorEventsMockFilterer
}

type OffchainAggregatorEventsMockCaller struct {
	contract *bind.BoundContract
}

type OffchainAggregatorEventsMockTransactor struct {
	contract *bind.BoundContract
}

type OffchainAggregatorEventsMockFilterer struct {
	contract *bind.BoundContract
}

type OffchainAggregatorEventsMockSession struct {
	Contract     *OffchainAggregatorEventsMock
	CallOpts     bind.CallOpts
	TransactOpts bind.TransactOpts
}

type OffchainAggregatorEventsMockCallerSession struct {
	Contract *OffchainAggregatorEventsMockCaller
	CallOpts bind.CallOpts
}

type OffchainAggregatorEventsMockTransactorSession struct {
	Contract     *OffchainAggregatorEventsMockTransactor
	TransactOpts bind.TransactOpts
}

type OffchainAggregatorEventsMockRaw struct {
	Contract *OffchainAggregatorEventsMock
}

type OffchainAggregatorEventsMockCallerRaw struct {
	Contract *OffchainAggregatorEventsMockCaller
}

type OffchainAggregatorEventsMockTransactorRaw struct {
	Contract *OffchainAggregatorEventsMockTransactor
}

func NewOffchainAggregatorEventsMock(address common.Address, backend bind.ContractBackend) (*OffchainAggregatorEventsMock, error) {
	abi, err := abi.JSON(strings.NewReader(OffchainAggregatorEventsMockABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindOffchainAggregatorEventsMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMock{address: address, abi: abi, OffchainAggregatorEventsMockCaller: OffchainAggregatorEventsMockCaller{contract: contract}, OffchainAggregatorEventsMockTransactor: OffchainAggregatorEventsMockTransactor{contract: contract}, OffchainAggregatorEventsMockFilterer: OffchainAggregatorEventsMockFilterer{contract: contract}}, nil
}

func NewOffchainAggregatorEventsMockCaller(address common.Address, caller bind.ContractCaller) (*OffchainAggregatorEventsMockCaller, error) {
	contract, err := bindOffchainAggregatorEventsMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockCaller{contract: contract}, nil
}

func NewOffchainAggregatorEventsMockTransactor(address common.Address, transactor bind.ContractTransactor) (*OffchainAggregatorEventsMockTransactor, error) {
	contract, err := bindOffchainAggregatorEventsMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockTransactor{contract: contract}, nil
}

func NewOffchainAggregatorEventsMockFilterer(address common.Address, filterer bind.ContractFilterer) (*OffchainAggregatorEventsMockFilterer, error) {
	contract, err := bindOffchainAggregatorEventsMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockFilterer{contract: contract}, nil
}

func bindOffchainAggregatorEventsMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OffchainAggregatorEventsMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffchainAggregatorEventsMock.Contract.OffchainAggregatorEventsMockCaller.contract.Call(opts, result, method, params...)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.OffchainAggregatorEventsMockTransactor.contract.Transfer(opts)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.OffchainAggregatorEventsMockTransactor.contract.Transact(opts, method, params...)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OffchainAggregatorEventsMock.Contract.contract.Call(opts, result, method, params...)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.contract.Transfer(opts)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.contract.Transact(opts, method, params...)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitAnswerUpdated(opts *bind.TransactOpts, current *big.Int, roundId *big.Int, updatedAt *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitAnswerUpdated", current, roundId, updatedAt)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitAnswerUpdated(current *big.Int, roundId *big.Int, updatedAt *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitAnswerUpdated(&_OffchainAggregatorEventsMock.TransactOpts, current, roundId, updatedAt)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitAnswerUpdated(current *big.Int, roundId *big.Int, updatedAt *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitAnswerUpdated(&_OffchainAggregatorEventsMock.TransactOpts, current, roundId, updatedAt)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitBillingAccessControllerSet(opts *bind.TransactOpts, old common.Address, current common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitBillingAccessControllerSet", old, current)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitBillingAccessControllerSet(old common.Address, current common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitBillingAccessControllerSet(&_OffchainAggregatorEventsMock.TransactOpts, old, current)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitBillingAccessControllerSet(old common.Address, current common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitBillingAccessControllerSet(&_OffchainAggregatorEventsMock.TransactOpts, old, current)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitBillingSet(opts *bind.TransactOpts, maximumGasPrice uint32, reasonableGasPrice uint32, microLinkPerEth uint32, linkGweiPerObservation uint32, linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitBillingSet", maximumGasPrice, reasonableGasPrice, microLinkPerEth, linkGweiPerObservation, linkGweiPerTransmission)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitBillingSet(maximumGasPrice uint32, reasonableGasPrice uint32, microLinkPerEth uint32, linkGweiPerObservation uint32, linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitBillingSet(&_OffchainAggregatorEventsMock.TransactOpts, maximumGasPrice, reasonableGasPrice, microLinkPerEth, linkGweiPerObservation, linkGweiPerTransmission)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitBillingSet(maximumGasPrice uint32, reasonableGasPrice uint32, microLinkPerEth uint32, linkGweiPerObservation uint32, linkGweiPerTransmission uint32) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitBillingSet(&_OffchainAggregatorEventsMock.TransactOpts, maximumGasPrice, reasonableGasPrice, microLinkPerEth, linkGweiPerObservation, linkGweiPerTransmission)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitConfigSet(opts *bind.TransactOpts, previousConfigBlockNumber uint32, configCount uint64, signers []common.Address, transmitters []common.Address, threshold uint8, encodedConfigVersion uint64, encoded []byte) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitConfigSet", previousConfigBlockNumber, configCount, signers, transmitters, threshold, encodedConfigVersion, encoded)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitConfigSet(previousConfigBlockNumber uint32, configCount uint64, signers []common.Address, transmitters []common.Address, threshold uint8, encodedConfigVersion uint64, encoded []byte) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitConfigSet(&_OffchainAggregatorEventsMock.TransactOpts, previousConfigBlockNumber, configCount, signers, transmitters, threshold, encodedConfigVersion, encoded)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitConfigSet(previousConfigBlockNumber uint32, configCount uint64, signers []common.Address, transmitters []common.Address, threshold uint8, encodedConfigVersion uint64, encoded []byte) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitConfigSet(&_OffchainAggregatorEventsMock.TransactOpts, previousConfigBlockNumber, configCount, signers, transmitters, threshold, encodedConfigVersion, encoded)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitLinkTokenSet(opts *bind.TransactOpts, _oldLinkToken common.Address, _newLinkToken common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitLinkTokenSet", _oldLinkToken, _newLinkToken)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitLinkTokenSet(_oldLinkToken common.Address, _newLinkToken common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitLinkTokenSet(&_OffchainAggregatorEventsMock.TransactOpts, _oldLinkToken, _newLinkToken)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitLinkTokenSet(_oldLinkToken common.Address, _newLinkToken common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitLinkTokenSet(&_OffchainAggregatorEventsMock.TransactOpts, _oldLinkToken, _newLinkToken)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitNewRound(opts *bind.TransactOpts, roundId *big.Int, startedBy common.Address, startedAt *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitNewRound", roundId, startedBy, startedAt)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitNewRound(roundId *big.Int, startedBy common.Address, startedAt *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitNewRound(&_OffchainAggregatorEventsMock.TransactOpts, roundId, startedBy, startedAt)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitNewRound(roundId *big.Int, startedBy common.Address, startedAt *big.Int) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitNewRound(&_OffchainAggregatorEventsMock.TransactOpts, roundId, startedBy, startedAt)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitNewTransmission(opts *bind.TransactOpts, aggregatorRoundId uint32, answer *big.Int, transmitter common.Address, observations []*big.Int, observers []byte, rawReportContext [32]byte) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitNewTransmission", aggregatorRoundId, answer, transmitter, observations, observers, rawReportContext)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitNewTransmission(aggregatorRoundId uint32, answer *big.Int, transmitter common.Address, observations []*big.Int, observers []byte, rawReportContext [32]byte) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitNewTransmission(&_OffchainAggregatorEventsMock.TransactOpts, aggregatorRoundId, answer, transmitter, observations, observers, rawReportContext)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitNewTransmission(aggregatorRoundId uint32, answer *big.Int, transmitter common.Address, observations []*big.Int, observers []byte, rawReportContext [32]byte) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitNewTransmission(&_OffchainAggregatorEventsMock.TransactOpts, aggregatorRoundId, answer, transmitter, observations, observers, rawReportContext)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitOraclePaid(opts *bind.TransactOpts, transmitter common.Address, payee common.Address, amount *big.Int, linkToken common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitOraclePaid", transmitter, payee, amount, linkToken)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitOraclePaid(transmitter common.Address, payee common.Address, amount *big.Int, linkToken common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitOraclePaid(&_OffchainAggregatorEventsMock.TransactOpts, transmitter, payee, amount, linkToken)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitOraclePaid(transmitter common.Address, payee common.Address, amount *big.Int, linkToken common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitOraclePaid(&_OffchainAggregatorEventsMock.TransactOpts, transmitter, payee, amount, linkToken)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitOwnershipTransferRequested", from, to)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitOwnershipTransferRequested(&_OffchainAggregatorEventsMock.TransactOpts, from, to)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitOwnershipTransferRequested(from common.Address, to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitOwnershipTransferRequested(&_OffchainAggregatorEventsMock.TransactOpts, from, to)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitOwnershipTransferred", from, to)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitOwnershipTransferred(&_OffchainAggregatorEventsMock.TransactOpts, from, to)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitOwnershipTransferred(from common.Address, to common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitOwnershipTransferred(&_OffchainAggregatorEventsMock.TransactOpts, from, to)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitPayeeshipTransferRequested(opts *bind.TransactOpts, transmitter common.Address, current common.Address, proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitPayeeshipTransferRequested", transmitter, current, proposed)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitPayeeshipTransferRequested(transmitter common.Address, current common.Address, proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitPayeeshipTransferRequested(&_OffchainAggregatorEventsMock.TransactOpts, transmitter, current, proposed)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitPayeeshipTransferRequested(transmitter common.Address, current common.Address, proposed common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitPayeeshipTransferRequested(&_OffchainAggregatorEventsMock.TransactOpts, transmitter, current, proposed)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitPayeeshipTransferred(opts *bind.TransactOpts, transmitter common.Address, previous common.Address, current common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitPayeeshipTransferred", transmitter, previous, current)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitPayeeshipTransferred(transmitter common.Address, previous common.Address, current common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitPayeeshipTransferred(&_OffchainAggregatorEventsMock.TransactOpts, transmitter, previous, current)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitPayeeshipTransferred(transmitter common.Address, previous common.Address, current common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitPayeeshipTransferred(&_OffchainAggregatorEventsMock.TransactOpts, transmitter, previous, current)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitRequesterAccessControllerSet(opts *bind.TransactOpts, old common.Address, current common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitRequesterAccessControllerSet", old, current)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitRequesterAccessControllerSet(old common.Address, current common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitRequesterAccessControllerSet(&_OffchainAggregatorEventsMock.TransactOpts, old, current)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitRequesterAccessControllerSet(old common.Address, current common.Address) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitRequesterAccessControllerSet(&_OffchainAggregatorEventsMock.TransactOpts, old, current)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitRoundRequested(opts *bind.TransactOpts, requester common.Address, configDigest [16]byte, epoch uint32, round uint8) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitRoundRequested", requester, configDigest, epoch, round)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitRoundRequested(requester common.Address, configDigest [16]byte, epoch uint32, round uint8) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitRoundRequested(&_OffchainAggregatorEventsMock.TransactOpts, requester, configDigest, epoch, round)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitRoundRequested(requester common.Address, configDigest [16]byte, epoch uint32, round uint8) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitRoundRequested(&_OffchainAggregatorEventsMock.TransactOpts, requester, configDigest, epoch, round)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactor) EmitValidatorConfigSet(opts *bind.TransactOpts, previousValidator common.Address, previousGasLimit uint32, currentValidator common.Address, currentGasLimit uint32) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.contract.Transact(opts, "emitValidatorConfigSet", previousValidator, previousGasLimit, currentValidator, currentGasLimit)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockSession) EmitValidatorConfigSet(previousValidator common.Address, previousGasLimit uint32, currentValidator common.Address, currentGasLimit uint32) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitValidatorConfigSet(&_OffchainAggregatorEventsMock.TransactOpts, previousValidator, previousGasLimit, currentValidator, currentGasLimit)
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockTransactorSession) EmitValidatorConfigSet(previousValidator common.Address, previousGasLimit uint32, currentValidator common.Address, currentGasLimit uint32) (*types.Transaction, error) {
	return _OffchainAggregatorEventsMock.Contract.EmitValidatorConfigSet(&_OffchainAggregatorEventsMock.TransactOpts, previousValidator, previousGasLimit, currentValidator, currentGasLimit)
}

type OffchainAggregatorEventsMockAnswerUpdatedIterator struct {
	Event *OffchainAggregatorEventsMockAnswerUpdated

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockAnswerUpdatedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockAnswerUpdated)
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
		it.Event = new(OffchainAggregatorEventsMockAnswerUpdated)
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

func (it *OffchainAggregatorEventsMockAnswerUpdatedIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockAnswerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockAnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
	Raw       types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*OffchainAggregatorEventsMockAnswerUpdatedIterator, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockAnswerUpdatedIterator{contract: _OffchainAggregatorEventsMock.contract, event: "AnswerUpdated", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error) {

	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "AnswerUpdated", currentRule, roundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockAnswerUpdated)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseAnswerUpdated(log types.Log) (*OffchainAggregatorEventsMockAnswerUpdated, error) {
	event := new(OffchainAggregatorEventsMockAnswerUpdated)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "AnswerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockBillingAccessControllerSetIterator struct {
	Event *OffchainAggregatorEventsMockBillingAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockBillingAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockBillingAccessControllerSet)
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
		it.Event = new(OffchainAggregatorEventsMockBillingAccessControllerSet)
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

func (it *OffchainAggregatorEventsMockBillingAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockBillingAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockBillingAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*OffchainAggregatorEventsMockBillingAccessControllerSetIterator, error) {

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockBillingAccessControllerSetIterator{contract: _OffchainAggregatorEventsMock.contract, event: "BillingAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockBillingAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "BillingAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockBillingAccessControllerSet)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseBillingAccessControllerSet(log types.Log) (*OffchainAggregatorEventsMockBillingAccessControllerSet, error) {
	event := new(OffchainAggregatorEventsMockBillingAccessControllerSet)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "BillingAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockBillingSetIterator struct {
	Event *OffchainAggregatorEventsMockBillingSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockBillingSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockBillingSet)
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
		it.Event = new(OffchainAggregatorEventsMockBillingSet)
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

func (it *OffchainAggregatorEventsMockBillingSetIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockBillingSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockBillingSet struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Raw                     types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterBillingSet(opts *bind.FilterOpts) (*OffchainAggregatorEventsMockBillingSetIterator, error) {

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockBillingSetIterator{contract: _OffchainAggregatorEventsMock.contract, event: "BillingSet", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchBillingSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockBillingSet) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "BillingSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockBillingSet)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "BillingSet", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseBillingSet(log types.Log) (*OffchainAggregatorEventsMockBillingSet, error) {
	event := new(OffchainAggregatorEventsMockBillingSet)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "BillingSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockConfigSetIterator struct {
	Event *OffchainAggregatorEventsMockConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockConfigSet)
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
		it.Event = new(OffchainAggregatorEventsMockConfigSet)
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

func (it *OffchainAggregatorEventsMockConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockConfigSet struct {
	PreviousConfigBlockNumber uint32
	ConfigCount               uint64
	Signers                   []common.Address
	Transmitters              []common.Address
	Threshold                 uint8
	EncodedConfigVersion      uint64
	Encoded                   []byte
	Raw                       types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterConfigSet(opts *bind.FilterOpts) (*OffchainAggregatorEventsMockConfigSetIterator, error) {

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockConfigSetIterator{contract: _OffchainAggregatorEventsMock.contract, event: "ConfigSet", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockConfigSet) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "ConfigSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockConfigSet)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "ConfigSet", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseConfigSet(log types.Log) (*OffchainAggregatorEventsMockConfigSet, error) {
	event := new(OffchainAggregatorEventsMockConfigSet)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "ConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockLinkTokenSetIterator struct {
	Event *OffchainAggregatorEventsMockLinkTokenSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockLinkTokenSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockLinkTokenSet)
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
		it.Event = new(OffchainAggregatorEventsMockLinkTokenSet)
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

func (it *OffchainAggregatorEventsMockLinkTokenSetIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockLinkTokenSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockLinkTokenSet struct {
	OldLinkToken common.Address
	NewLinkToken common.Address
	Raw          types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterLinkTokenSet(opts *bind.FilterOpts, _oldLinkToken []common.Address, _newLinkToken []common.Address) (*OffchainAggregatorEventsMockLinkTokenSetIterator, error) {

	var _oldLinkTokenRule []interface{}
	for _, _oldLinkTokenItem := range _oldLinkToken {
		_oldLinkTokenRule = append(_oldLinkTokenRule, _oldLinkTokenItem)
	}
	var _newLinkTokenRule []interface{}
	for _, _newLinkTokenItem := range _newLinkToken {
		_newLinkTokenRule = append(_newLinkTokenRule, _newLinkTokenItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "LinkTokenSet", _oldLinkTokenRule, _newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockLinkTokenSetIterator{contract: _OffchainAggregatorEventsMock.contract, event: "LinkTokenSet", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockLinkTokenSet, _oldLinkToken []common.Address, _newLinkToken []common.Address) (event.Subscription, error) {

	var _oldLinkTokenRule []interface{}
	for _, _oldLinkTokenItem := range _oldLinkToken {
		_oldLinkTokenRule = append(_oldLinkTokenRule, _oldLinkTokenItem)
	}
	var _newLinkTokenRule []interface{}
	for _, _newLinkTokenItem := range _newLinkToken {
		_newLinkTokenRule = append(_newLinkTokenRule, _newLinkTokenItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "LinkTokenSet", _oldLinkTokenRule, _newLinkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockLinkTokenSet)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseLinkTokenSet(log types.Log) (*OffchainAggregatorEventsMockLinkTokenSet, error) {
	event := new(OffchainAggregatorEventsMockLinkTokenSet)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "LinkTokenSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockNewRoundIterator struct {
	Event *OffchainAggregatorEventsMockNewRound

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockNewRoundIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockNewRound)
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
		it.Event = new(OffchainAggregatorEventsMockNewRound)
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

func (it *OffchainAggregatorEventsMockNewRoundIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockNewRoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockNewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
	Raw       types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*OffchainAggregatorEventsMockNewRoundIterator, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockNewRoundIterator{contract: _OffchainAggregatorEventsMock.contract, event: "NewRound", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchNewRound(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error) {

	var roundIdRule []interface{}
	for _, roundIdItem := range roundId {
		roundIdRule = append(roundIdRule, roundIdItem)
	}
	var startedByRule []interface{}
	for _, startedByItem := range startedBy {
		startedByRule = append(startedByRule, startedByItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "NewRound", roundIdRule, startedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockNewRound)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "NewRound", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseNewRound(log types.Log) (*OffchainAggregatorEventsMockNewRound, error) {
	event := new(OffchainAggregatorEventsMockNewRound)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "NewRound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockNewTransmissionIterator struct {
	Event *OffchainAggregatorEventsMockNewTransmission

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockNewTransmissionIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockNewTransmission)
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
		it.Event = new(OffchainAggregatorEventsMockNewTransmission)
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

func (it *OffchainAggregatorEventsMockNewTransmissionIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockNewTransmissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockNewTransmission struct {
	AggregatorRoundId uint32
	Answer            *big.Int
	Transmitter       common.Address
	Observations      []*big.Int
	Observers         []byte
	RawReportContext  [32]byte
	Raw               types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*OffchainAggregatorEventsMockNewTransmissionIterator, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockNewTransmissionIterator{contract: _OffchainAggregatorEventsMock.contract, event: "NewTransmission", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error) {

	var aggregatorRoundIdRule []interface{}
	for _, aggregatorRoundIdItem := range aggregatorRoundId {
		aggregatorRoundIdRule = append(aggregatorRoundIdRule, aggregatorRoundIdItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "NewTransmission", aggregatorRoundIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockNewTransmission)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "NewTransmission", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseNewTransmission(log types.Log) (*OffchainAggregatorEventsMockNewTransmission, error) {
	event := new(OffchainAggregatorEventsMockNewTransmission)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "NewTransmission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockOraclePaidIterator struct {
	Event *OffchainAggregatorEventsMockOraclePaid

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockOraclePaidIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockOraclePaid)
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
		it.Event = new(OffchainAggregatorEventsMockOraclePaid)
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

func (it *OffchainAggregatorEventsMockOraclePaidIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockOraclePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockOraclePaid struct {
	Transmitter common.Address
	Payee       common.Address
	Amount      *big.Int
	LinkToken   common.Address
	Raw         types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*OffchainAggregatorEventsMockOraclePaidIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	var linkTokenRule []interface{}
	for _, linkTokenItem := range linkToken {
		linkTokenRule = append(linkTokenRule, linkTokenItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockOraclePaidIterator{contract: _OffchainAggregatorEventsMock.contract, event: "OraclePaid", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var payeeRule []interface{}
	for _, payeeItem := range payee {
		payeeRule = append(payeeRule, payeeItem)
	}

	var linkTokenRule []interface{}
	for _, linkTokenItem := range linkToken {
		linkTokenRule = append(linkTokenRule, linkTokenItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "OraclePaid", transmitterRule, payeeRule, linkTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockOraclePaid)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "OraclePaid", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseOraclePaid(log types.Log) (*OffchainAggregatorEventsMockOraclePaid, error) {
	event := new(OffchainAggregatorEventsMockOraclePaid)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "OraclePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockOwnershipTransferRequestedIterator struct {
	Event *OffchainAggregatorEventsMockOwnershipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockOwnershipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockOwnershipTransferRequested)
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
		it.Event = new(OffchainAggregatorEventsMockOwnershipTransferRequested)
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

func (it *OffchainAggregatorEventsMockOwnershipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockOwnershipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockOwnershipTransferRequested struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffchainAggregatorEventsMockOwnershipTransferRequestedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockOwnershipTransferRequestedIterator{contract: _OffchainAggregatorEventsMock.contract, event: "OwnershipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "OwnershipTransferRequested", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockOwnershipTransferRequested)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseOwnershipTransferRequested(log types.Log) (*OffchainAggregatorEventsMockOwnershipTransferRequested, error) {
	event := new(OffchainAggregatorEventsMockOwnershipTransferRequested)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "OwnershipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockOwnershipTransferredIterator struct {
	Event *OffchainAggregatorEventsMockOwnershipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockOwnershipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockOwnershipTransferred)
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
		it.Event = new(OffchainAggregatorEventsMockOwnershipTransferred)
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

func (it *OffchainAggregatorEventsMockOwnershipTransferredIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffchainAggregatorEventsMockOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockOwnershipTransferredIterator{contract: _OffchainAggregatorEventsMock.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockOwnershipTransferred)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseOwnershipTransferred(log types.Log) (*OffchainAggregatorEventsMockOwnershipTransferred, error) {
	event := new(OffchainAggregatorEventsMockOwnershipTransferred)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockPayeeshipTransferRequestedIterator struct {
	Event *OffchainAggregatorEventsMockPayeeshipTransferRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockPayeeshipTransferRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockPayeeshipTransferRequested)
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
		it.Event = new(OffchainAggregatorEventsMockPayeeshipTransferRequested)
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

func (it *OffchainAggregatorEventsMockPayeeshipTransferRequestedIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockPayeeshipTransferRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockPayeeshipTransferRequested struct {
	Transmitter common.Address
	Current     common.Address
	Proposed    common.Address
	Raw         types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*OffchainAggregatorEventsMockPayeeshipTransferRequestedIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockPayeeshipTransferRequestedIterator{contract: _OffchainAggregatorEventsMock.contract, event: "PayeeshipTransferRequested", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "PayeeshipTransferRequested", transmitterRule, currentRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockPayeeshipTransferRequested)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParsePayeeshipTransferRequested(log types.Log) (*OffchainAggregatorEventsMockPayeeshipTransferRequested, error) {
	event := new(OffchainAggregatorEventsMockPayeeshipTransferRequested)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "PayeeshipTransferRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockPayeeshipTransferredIterator struct {
	Event *OffchainAggregatorEventsMockPayeeshipTransferred

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockPayeeshipTransferredIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockPayeeshipTransferred)
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
		it.Event = new(OffchainAggregatorEventsMockPayeeshipTransferred)
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

func (it *OffchainAggregatorEventsMockPayeeshipTransferredIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockPayeeshipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockPayeeshipTransferred struct {
	Transmitter common.Address
	Previous    common.Address
	Current     common.Address
	Raw         types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*OffchainAggregatorEventsMockPayeeshipTransferredIterator, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockPayeeshipTransferredIterator{contract: _OffchainAggregatorEventsMock.contract, event: "PayeeshipTransferred", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error) {

	var transmitterRule []interface{}
	for _, transmitterItem := range transmitter {
		transmitterRule = append(transmitterRule, transmitterItem)
	}
	var previousRule []interface{}
	for _, previousItem := range previous {
		previousRule = append(previousRule, previousItem)
	}
	var currentRule []interface{}
	for _, currentItem := range current {
		currentRule = append(currentRule, currentItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "PayeeshipTransferred", transmitterRule, previousRule, currentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockPayeeshipTransferred)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParsePayeeshipTransferred(log types.Log) (*OffchainAggregatorEventsMockPayeeshipTransferred, error) {
	event := new(OffchainAggregatorEventsMockPayeeshipTransferred)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "PayeeshipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockRequesterAccessControllerSetIterator struct {
	Event *OffchainAggregatorEventsMockRequesterAccessControllerSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockRequesterAccessControllerSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockRequesterAccessControllerSet)
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
		it.Event = new(OffchainAggregatorEventsMockRequesterAccessControllerSet)
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

func (it *OffchainAggregatorEventsMockRequesterAccessControllerSetIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockRequesterAccessControllerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockRequesterAccessControllerSet struct {
	Old     common.Address
	Current common.Address
	Raw     types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*OffchainAggregatorEventsMockRequesterAccessControllerSetIterator, error) {

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockRequesterAccessControllerSetIterator{contract: _OffchainAggregatorEventsMock.contract, event: "RequesterAccessControllerSet", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockRequesterAccessControllerSet) (event.Subscription, error) {

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "RequesterAccessControllerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockRequesterAccessControllerSet)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseRequesterAccessControllerSet(log types.Log) (*OffchainAggregatorEventsMockRequesterAccessControllerSet, error) {
	event := new(OffchainAggregatorEventsMockRequesterAccessControllerSet)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "RequesterAccessControllerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockRoundRequestedIterator struct {
	Event *OffchainAggregatorEventsMockRoundRequested

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockRoundRequestedIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockRoundRequested)
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
		it.Event = new(OffchainAggregatorEventsMockRoundRequested)
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

func (it *OffchainAggregatorEventsMockRoundRequestedIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockRoundRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockRoundRequested struct {
	Requester    common.Address
	ConfigDigest [16]byte
	Epoch        uint32
	Round        uint8
	Raw          types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*OffchainAggregatorEventsMockRoundRequestedIterator, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockRoundRequestedIterator{contract: _OffchainAggregatorEventsMock.contract, event: "RoundRequested", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockRoundRequested, requester []common.Address) (event.Subscription, error) {

	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "RoundRequested", requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockRoundRequested)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "RoundRequested", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseRoundRequested(log types.Log) (*OffchainAggregatorEventsMockRoundRequested, error) {
	event := new(OffchainAggregatorEventsMockRoundRequested)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "RoundRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type OffchainAggregatorEventsMockValidatorConfigSetIterator struct {
	Event *OffchainAggregatorEventsMockValidatorConfigSet

	contract *bind.BoundContract
	event    string

	logs chan types.Log
	sub  ethereum.Subscription
	done bool
	fail error
}

func (it *OffchainAggregatorEventsMockValidatorConfigSetIterator) Next() bool {

	if it.fail != nil {
		return false
	}

	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OffchainAggregatorEventsMockValidatorConfigSet)
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
		it.Event = new(OffchainAggregatorEventsMockValidatorConfigSet)
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

func (it *OffchainAggregatorEventsMockValidatorConfigSetIterator) Error() error {
	return it.fail
}

func (it *OffchainAggregatorEventsMockValidatorConfigSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type OffchainAggregatorEventsMockValidatorConfigSet struct {
	PreviousValidator common.Address
	PreviousGasLimit  uint32
	CurrentValidator  common.Address
	CurrentGasLimit   uint32
	Raw               types.Log
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*OffchainAggregatorEventsMockValidatorConfigSetIterator, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.FilterLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return &OffchainAggregatorEventsMockValidatorConfigSetIterator{contract: _OffchainAggregatorEventsMock.contract, event: "ValidatorConfigSet", logs: logs, sub: sub}, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error) {

	var previousValidatorRule []interface{}
	for _, previousValidatorItem := range previousValidator {
		previousValidatorRule = append(previousValidatorRule, previousValidatorItem)
	}

	var currentValidatorRule []interface{}
	for _, currentValidatorItem := range currentValidator {
		currentValidatorRule = append(currentValidatorRule, currentValidatorItem)
	}

	logs, sub, err := _OffchainAggregatorEventsMock.contract.WatchLogs(opts, "ValidatorConfigSet", previousValidatorRule, currentValidatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:

				event := new(OffchainAggregatorEventsMockValidatorConfigSet)
				if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
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

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMockFilterer) ParseValidatorConfigSet(log types.Log) (*OffchainAggregatorEventsMockValidatorConfigSet, error) {
	event := new(OffchainAggregatorEventsMockValidatorConfigSet)
	if err := _OffchainAggregatorEventsMock.contract.UnpackLog(event, "ValidatorConfigSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMock) ParseLog(log types.Log) (generated.AbigenLog, error) {
	switch log.Topics[0] {
	case _OffchainAggregatorEventsMock.abi.Events["AnswerUpdated"].ID:
		return _OffchainAggregatorEventsMock.ParseAnswerUpdated(log)
	case _OffchainAggregatorEventsMock.abi.Events["BillingAccessControllerSet"].ID:
		return _OffchainAggregatorEventsMock.ParseBillingAccessControllerSet(log)
	case _OffchainAggregatorEventsMock.abi.Events["BillingSet"].ID:
		return _OffchainAggregatorEventsMock.ParseBillingSet(log)
	case _OffchainAggregatorEventsMock.abi.Events["ConfigSet"].ID:
		return _OffchainAggregatorEventsMock.ParseConfigSet(log)
	case _OffchainAggregatorEventsMock.abi.Events["LinkTokenSet"].ID:
		return _OffchainAggregatorEventsMock.ParseLinkTokenSet(log)
	case _OffchainAggregatorEventsMock.abi.Events["NewRound"].ID:
		return _OffchainAggregatorEventsMock.ParseNewRound(log)
	case _OffchainAggregatorEventsMock.abi.Events["NewTransmission"].ID:
		return _OffchainAggregatorEventsMock.ParseNewTransmission(log)
	case _OffchainAggregatorEventsMock.abi.Events["OraclePaid"].ID:
		return _OffchainAggregatorEventsMock.ParseOraclePaid(log)
	case _OffchainAggregatorEventsMock.abi.Events["OwnershipTransferRequested"].ID:
		return _OffchainAggregatorEventsMock.ParseOwnershipTransferRequested(log)
	case _OffchainAggregatorEventsMock.abi.Events["OwnershipTransferred"].ID:
		return _OffchainAggregatorEventsMock.ParseOwnershipTransferred(log)
	case _OffchainAggregatorEventsMock.abi.Events["PayeeshipTransferRequested"].ID:
		return _OffchainAggregatorEventsMock.ParsePayeeshipTransferRequested(log)
	case _OffchainAggregatorEventsMock.abi.Events["PayeeshipTransferred"].ID:
		return _OffchainAggregatorEventsMock.ParsePayeeshipTransferred(log)
	case _OffchainAggregatorEventsMock.abi.Events["RequesterAccessControllerSet"].ID:
		return _OffchainAggregatorEventsMock.ParseRequesterAccessControllerSet(log)
	case _OffchainAggregatorEventsMock.abi.Events["RoundRequested"].ID:
		return _OffchainAggregatorEventsMock.ParseRoundRequested(log)
	case _OffchainAggregatorEventsMock.abi.Events["ValidatorConfigSet"].ID:
		return _OffchainAggregatorEventsMock.ParseValidatorConfigSet(log)

	default:
		return nil, fmt.Errorf("abigen wrapper received unknown log topic: %v", log.Topics[0])
	}
}

func (OffchainAggregatorEventsMockAnswerUpdated) Topic() common.Hash {
	return common.HexToHash("0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f")
}

func (OffchainAggregatorEventsMockBillingAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x793cb73064f3c8cde7e187ae515511e6e56d1ee89bf08b82fa60fb70f8d48912")
}

func (OffchainAggregatorEventsMockBillingSet) Topic() common.Hash {
	return common.HexToHash("0xd0d9486a2c673e2a4b57fc82e4c8a556b3e2b82dd5db07e2c04a920ca0f469b6")
}

func (OffchainAggregatorEventsMockConfigSet) Topic() common.Hash {
	return common.HexToHash("0x25d719d88a4512dd76c7442b910a83360845505894eb444ef299409e180f8fb9")
}

func (OffchainAggregatorEventsMockLinkTokenSet) Topic() common.Hash {
	return common.HexToHash("0x4966a50c93f855342ccf6c5c0d358b85b91335b2acedc7da0932f691f351711a")
}

func (OffchainAggregatorEventsMockNewRound) Topic() common.Hash {
	return common.HexToHash("0x0109fc6f55cf40689f02fbaad7af7fe7bbac8a3d2186600afc7d3e10cac60271")
}

func (OffchainAggregatorEventsMockNewTransmission) Topic() common.Hash {
	return common.HexToHash("0xf6a97944f31ea060dfde0566e4167c1a1082551e64b60ecb14d599a9d023d451")
}

func (OffchainAggregatorEventsMockOraclePaid) Topic() common.Hash {
	return common.HexToHash("0xd0b1dac935d85bd54cf0a33b0d41d39f8cf53a968465fc7ea2377526b8ac712c")
}

func (OffchainAggregatorEventsMockOwnershipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0xed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae1278")
}

func (OffchainAggregatorEventsMockOwnershipTransferred) Topic() common.Hash {
	return common.HexToHash("0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0")
}

func (OffchainAggregatorEventsMockPayeeshipTransferRequested) Topic() common.Hash {
	return common.HexToHash("0x84f7c7c80bb8ed2279b4aab5f61cd05e6374073d38f46d7f32de8c30e9e38367")
}

func (OffchainAggregatorEventsMockPayeeshipTransferred) Topic() common.Hash {
	return common.HexToHash("0x78af32efdcad432315431e9b03d27e6cd98fb79c405fdc5af7c1714d9c0f75b3")
}

func (OffchainAggregatorEventsMockRequesterAccessControllerSet) Topic() common.Hash {
	return common.HexToHash("0x27b89aede8b560578baaa25ee5ce3852c5eecad1e114b941bbd89e1eb4bae634")
}

func (OffchainAggregatorEventsMockRoundRequested) Topic() common.Hash {
	return common.HexToHash("0x3ea16a923ff4b1df6526e854c9e3a995c43385d70e73359e10623c74f0b52037")
}

func (OffchainAggregatorEventsMockValidatorConfigSet) Topic() common.Hash {
	return common.HexToHash("0xb04e3a37abe9c0fcdfebdeae019a8e2b12ddf53f5d55ffb0caccc1bedaca1541")
}

func (_OffchainAggregatorEventsMock *OffchainAggregatorEventsMock) Address() common.Address {
	return _OffchainAggregatorEventsMock.address
}

type OffchainAggregatorEventsMockInterface interface {
	EmitAnswerUpdated(opts *bind.TransactOpts, current *big.Int, roundId *big.Int, updatedAt *big.Int) (*types.Transaction, error)

	EmitBillingAccessControllerSet(opts *bind.TransactOpts, old common.Address, current common.Address) (*types.Transaction, error)

	EmitBillingSet(opts *bind.TransactOpts, maximumGasPrice uint32, reasonableGasPrice uint32, microLinkPerEth uint32, linkGweiPerObservation uint32, linkGweiPerTransmission uint32) (*types.Transaction, error)

	EmitConfigSet(opts *bind.TransactOpts, previousConfigBlockNumber uint32, configCount uint64, signers []common.Address, transmitters []common.Address, threshold uint8, encodedConfigVersion uint64, encoded []byte) (*types.Transaction, error)

	EmitLinkTokenSet(opts *bind.TransactOpts, _oldLinkToken common.Address, _newLinkToken common.Address) (*types.Transaction, error)

	EmitNewRound(opts *bind.TransactOpts, roundId *big.Int, startedBy common.Address, startedAt *big.Int) (*types.Transaction, error)

	EmitNewTransmission(opts *bind.TransactOpts, aggregatorRoundId uint32, answer *big.Int, transmitter common.Address, observations []*big.Int, observers []byte, rawReportContext [32]byte) (*types.Transaction, error)

	EmitOraclePaid(opts *bind.TransactOpts, transmitter common.Address, payee common.Address, amount *big.Int, linkToken common.Address) (*types.Transaction, error)

	EmitOwnershipTransferRequested(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitOwnershipTransferred(opts *bind.TransactOpts, from common.Address, to common.Address) (*types.Transaction, error)

	EmitPayeeshipTransferRequested(opts *bind.TransactOpts, transmitter common.Address, current common.Address, proposed common.Address) (*types.Transaction, error)

	EmitPayeeshipTransferred(opts *bind.TransactOpts, transmitter common.Address, previous common.Address, current common.Address) (*types.Transaction, error)

	EmitRequesterAccessControllerSet(opts *bind.TransactOpts, old common.Address, current common.Address) (*types.Transaction, error)

	EmitRoundRequested(opts *bind.TransactOpts, requester common.Address, configDigest [16]byte, epoch uint32, round uint8) (*types.Transaction, error)

	EmitValidatorConfigSet(opts *bind.TransactOpts, previousValidator common.Address, previousGasLimit uint32, currentValidator common.Address, currentGasLimit uint32) (*types.Transaction, error)

	FilterAnswerUpdated(opts *bind.FilterOpts, current []*big.Int, roundId []*big.Int) (*OffchainAggregatorEventsMockAnswerUpdatedIterator, error)

	WatchAnswerUpdated(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockAnswerUpdated, current []*big.Int, roundId []*big.Int) (event.Subscription, error)

	ParseAnswerUpdated(log types.Log) (*OffchainAggregatorEventsMockAnswerUpdated, error)

	FilterBillingAccessControllerSet(opts *bind.FilterOpts) (*OffchainAggregatorEventsMockBillingAccessControllerSetIterator, error)

	WatchBillingAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockBillingAccessControllerSet) (event.Subscription, error)

	ParseBillingAccessControllerSet(log types.Log) (*OffchainAggregatorEventsMockBillingAccessControllerSet, error)

	FilterBillingSet(opts *bind.FilterOpts) (*OffchainAggregatorEventsMockBillingSetIterator, error)

	WatchBillingSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockBillingSet) (event.Subscription, error)

	ParseBillingSet(log types.Log) (*OffchainAggregatorEventsMockBillingSet, error)

	FilterConfigSet(opts *bind.FilterOpts) (*OffchainAggregatorEventsMockConfigSetIterator, error)

	WatchConfigSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockConfigSet) (event.Subscription, error)

	ParseConfigSet(log types.Log) (*OffchainAggregatorEventsMockConfigSet, error)

	FilterLinkTokenSet(opts *bind.FilterOpts, _oldLinkToken []common.Address, _newLinkToken []common.Address) (*OffchainAggregatorEventsMockLinkTokenSetIterator, error)

	WatchLinkTokenSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockLinkTokenSet, _oldLinkToken []common.Address, _newLinkToken []common.Address) (event.Subscription, error)

	ParseLinkTokenSet(log types.Log) (*OffchainAggregatorEventsMockLinkTokenSet, error)

	FilterNewRound(opts *bind.FilterOpts, roundId []*big.Int, startedBy []common.Address) (*OffchainAggregatorEventsMockNewRoundIterator, error)

	WatchNewRound(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockNewRound, roundId []*big.Int, startedBy []common.Address) (event.Subscription, error)

	ParseNewRound(log types.Log) (*OffchainAggregatorEventsMockNewRound, error)

	FilterNewTransmission(opts *bind.FilterOpts, aggregatorRoundId []uint32) (*OffchainAggregatorEventsMockNewTransmissionIterator, error)

	WatchNewTransmission(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockNewTransmission, aggregatorRoundId []uint32) (event.Subscription, error)

	ParseNewTransmission(log types.Log) (*OffchainAggregatorEventsMockNewTransmission, error)

	FilterOraclePaid(opts *bind.FilterOpts, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (*OffchainAggregatorEventsMockOraclePaidIterator, error)

	WatchOraclePaid(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockOraclePaid, transmitter []common.Address, payee []common.Address, linkToken []common.Address) (event.Subscription, error)

	ParseOraclePaid(log types.Log) (*OffchainAggregatorEventsMockOraclePaid, error)

	FilterOwnershipTransferRequested(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffchainAggregatorEventsMockOwnershipTransferRequestedIterator, error)

	WatchOwnershipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockOwnershipTransferRequested, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferRequested(log types.Log) (*OffchainAggregatorEventsMockOwnershipTransferRequested, error)

	FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OffchainAggregatorEventsMockOwnershipTransferredIterator, error)

	WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error)

	ParseOwnershipTransferred(log types.Log) (*OffchainAggregatorEventsMockOwnershipTransferred, error)

	FilterPayeeshipTransferRequested(opts *bind.FilterOpts, transmitter []common.Address, current []common.Address, proposed []common.Address) (*OffchainAggregatorEventsMockPayeeshipTransferRequestedIterator, error)

	WatchPayeeshipTransferRequested(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockPayeeshipTransferRequested, transmitter []common.Address, current []common.Address, proposed []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferRequested(log types.Log) (*OffchainAggregatorEventsMockPayeeshipTransferRequested, error)

	FilterPayeeshipTransferred(opts *bind.FilterOpts, transmitter []common.Address, previous []common.Address, current []common.Address) (*OffchainAggregatorEventsMockPayeeshipTransferredIterator, error)

	WatchPayeeshipTransferred(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockPayeeshipTransferred, transmitter []common.Address, previous []common.Address, current []common.Address) (event.Subscription, error)

	ParsePayeeshipTransferred(log types.Log) (*OffchainAggregatorEventsMockPayeeshipTransferred, error)

	FilterRequesterAccessControllerSet(opts *bind.FilterOpts) (*OffchainAggregatorEventsMockRequesterAccessControllerSetIterator, error)

	WatchRequesterAccessControllerSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockRequesterAccessControllerSet) (event.Subscription, error)

	ParseRequesterAccessControllerSet(log types.Log) (*OffchainAggregatorEventsMockRequesterAccessControllerSet, error)

	FilterRoundRequested(opts *bind.FilterOpts, requester []common.Address) (*OffchainAggregatorEventsMockRoundRequestedIterator, error)

	WatchRoundRequested(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockRoundRequested, requester []common.Address) (event.Subscription, error)

	ParseRoundRequested(log types.Log) (*OffchainAggregatorEventsMockRoundRequested, error)

	FilterValidatorConfigSet(opts *bind.FilterOpts, previousValidator []common.Address, currentValidator []common.Address) (*OffchainAggregatorEventsMockValidatorConfigSetIterator, error)

	WatchValidatorConfigSet(opts *bind.WatchOpts, sink chan<- *OffchainAggregatorEventsMockValidatorConfigSet, previousValidator []common.Address, currentValidator []common.Address) (event.Subscription, error)

	ParseValidatorConfigSet(log types.Log) (*OffchainAggregatorEventsMockValidatorConfigSet, error)

	ParseLog(log types.Log) (generated.AbigenLog, error)

	Address() common.Address
}
