package transmission_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_consumer_interface_v08"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/vrf_coordinator_mock"

	"github.com/ethereum/go-ethereum/core/types"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/transmission/generated/entry_point"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/transmission/generated/greeter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/transmission/generated/paymaster_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/transmission/generated/sca_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/transmission/generated/smart_contract_account_factory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/transmission/generated/smart_contract_account_helper"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/transmission"
)

var (
	greeterABI    = evmtypes.MustGetABI(greeter_wrapper.GreeterABI)
	consumerABI   = evmtypes.MustGetABI(solidity_vrf_consumer_interface_v08.VRFConsumerABI)
	entrypointABI = evmtypes.MustGetABI(entry_point.EntryPointABI)
)

type EntryPointUniverse struct {
	holder1               *bind.TransactOpts
	holder1Key            ethkey.KeyV2
	holder2               *bind.TransactOpts
	backend               *backends.SimulatedBackend
	entryPointAddress     common.Address
	entryPoint            *entry_point.EntryPoint
	factoryAddress        common.Address
	helper                *smart_contract_account_helper.SmartContractAccountHelper
	greeterAddress        common.Address
	greeter               *greeter_wrapper.Greeter
	linkTokenAddress      common.Address
	linkToken             *link_token_interface.LinkToken
	linkEthFeedAddress    common.Address
	vrfCoordinatorAddress common.Address
	vrfCoordinator        *vrf_coordinator_mock.VRFCoordinatorMock
	vrfConsumerAddress    common.Address
}

func deployTransmissionUniverse(t *testing.T) *EntryPointUniverse {
	// Create a key for holder1 that we can use to sign
	holder1Key := cltest.MustGenerateRandomKey(t)
	t.Log("Holder key:", holder1Key.String())

	// Construct simulated blockchain environment.
	holder1Transactor, err := bind.NewKeyedTransactorWithChainID(holder1Key.ToEcdsaPrivKey(), testutils.SimulatedChainID)
	require.NoError(t, err)
	var (
		holder1 = holder1Transactor
		holder2 = testutils.MustNewSimTransactor(t)
	)
	genesisData := core.GenesisAlloc{
		holder1.From: {Balance: assets.Ether(1000).ToInt()},
		holder2.From: {Balance: assets.Ether(1000).ToInt()},
	}
	gasLimit := uint32(30e6)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	backend.Commit()

	// Setup all contracts and addresses used by tests.
	entryPointAddress, _, entryPoint, err := entry_point.DeployEntryPoint(holder1, backend)
	require.NoError(t, err)
	factoryAddress, _, _, _ := smart_contract_account_factory.DeploySmartContractAccountFactory(holder1, backend)
	require.NoError(t, err)
	_, _, helper, err := smart_contract_account_helper.DeploySmartContractAccountHelper(holder1, backend)
	require.NoError(t, err)
	greeterAddress, _, greeter, err := greeter_wrapper.DeployGreeter(holder1, backend)
	require.NoError(t, err)
	linkTokenAddress, _, linkToken, err := link_token_interface.DeployLinkToken(holder1, backend)
	require.NoError(t, err)
	linkEthFeedAddress, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
		holder1,
		backend,
		18,
		(*big.Int)(assets.GWei(5000000)), // .005 ETH
	)
	require.NoError(t, err)
	vrfCoordinatorAddress, _, vrfCoordinator, err := vrf_coordinator_mock.DeployVRFCoordinatorMock(holder1, backend, linkTokenAddress)
	require.NoError(t, err)
	vrfConsumerAddress, _, _, err := solidity_vrf_consumer_interface_v08.DeployVRFConsumer(holder1, backend, vrfCoordinatorAddress, linkTokenAddress)
	require.NoError(t, err)
	backend.Commit()

	return &EntryPointUniverse{
		holder1:               holder1,
		holder1Key:            holder1Key,
		holder2:               holder2,
		backend:               backend,
		entryPointAddress:     entryPointAddress,
		entryPoint:            entryPoint,
		factoryAddress:        factoryAddress,
		helper:                helper,
		greeterAddress:        greeterAddress,
		greeter:               greeter,
		linkTokenAddress:      linkTokenAddress,
		linkToken:             linkToken,
		linkEthFeedAddress:    linkEthFeedAddress,
		vrfCoordinatorAddress: vrfCoordinatorAddress,
		vrfCoordinator:        vrfCoordinator,
		vrfConsumerAddress:    vrfConsumerAddress,
	}
}

func Test4337Basic(t *testing.T) {
	// Deploy universe.
	universe := deployTransmissionUniverse(t)
	holder1 := universe.holder1
	holder2 := universe.holder2
	backend := universe.backend

	// Ensure no greeting is already set.
	initialGreeting, err := universe.greeter.GetGreeting(nil)
	require.NoError(t, err)
	require.Equal(t, "", initialGreeting)

	// Get the address at which the Smart Contract Account will be deployed.
	toDeployAddress, err := universe.helper.CalculateSmartContractAccountAddress(
		nil,
		holder1.From,
		universe.entryPointAddress,
		universe.factoryAddress,
	)
	require.NoError(t, err)
	t.Log("Smart Contract Account Address:", toDeployAddress)

	// Get the initialization code for the Smart Contract Account.
	fullInitializeCode, err := universe.helper.GetInitCode(nil, universe.factoryAddress, holder1.From, universe.entryPointAddress)
	require.NoError(t, err)
	t.Log("Full initialization code:", common.Bytes2Hex(fullInitializeCode))

	// Construct calldata for setGreeting.
	encodedGreetingCall, err := greeterABI.Pack("setGreeting", "bye")
	require.NoError(t, err)
	t.Log("Encoded greeting call:", common.Bytes2Hex(encodedGreetingCall))

	// Construct the calldata to be passed in the user operation.
	var (
		value    = big.NewInt(0)
		nonce    = big.NewInt(0)
		deadline = big.NewInt(1000)
	)
	fullEncoding, err := universe.helper.GetFullEndTxEncoding(nil, universe.greeterAddress, value, deadline, encodedGreetingCall)
	require.NoError(t, err)
	t.Log("Full user operation calldata:", common.Bytes2Hex(fullEncoding))

	// Construct and execute user operation.
	userOp := entry_point.UserOperation{
		Sender:               toDeployAddress,
		Nonce:                nonce,
		InitCode:             fullInitializeCode,
		CallData:             fullEncoding,
		CallGasLimit:         big.NewInt(10_000_000),
		VerificationGasLimit: big.NewInt(10_000_000),
		PreVerificationGas:   big.NewInt(10_000_000),
		MaxFeePerGas:         big.NewInt(100),
		MaxPriorityFeePerGas: big.NewInt(200),
		PaymasterAndData:     []byte(""),
		Signature:            []byte(""),
	}

	// Generate hash from user operation, sign it, and include it in the user operation.
	userOpHash, err := universe.entryPoint.GetUserOpHash(nil, userOp)
	require.NoError(t, err)
	fullHash, err := universe.helper.GetFullHashForSigning(nil, userOpHash, toDeployAddress)
	require.NoError(t, err)
	t.Log("Full hash for signing:", common.Bytes2Hex(fullHash[:]))
	sig, err := transmission.SignMessage(universe.holder1Key.ToEcdsaPrivKey(), fullHash[:])
	require.NoError(t, err)
	t.Log("Signature:", common.Bytes2Hex(sig))
	userOp.Signature = sig

	// Deposit to the SCA's account to pay for this transaction.
	holder1.Value = assets.Ether(10).ToInt()
	tx, err := universe.entryPoint.DepositTo(holder1, toDeployAddress)
	require.NoError(t, err)
	backend.Commit()
	_, err = bind.WaitMined(testutils.Context(t), backend, tx)
	require.NoError(t, err)
	holder1.Value = assets.Ether(0).ToInt()
	balance, err := universe.entryPoint.BalanceOf(nil, toDeployAddress)
	require.NoError(t, err)
	require.Equal(t, assets.Ether(10).ToInt(), balance)

	// Run handleOps from holder2's account, to demonstrate that any account can execute this signed user operation.
	tx, err = universe.entryPoint.HandleOps(holder2, []entry_point.UserOperation{userOp}, holder1.From)
	require.NoError(t, err)
	backend.Commit()
	_, err = bind.WaitMined(testutils.Context(t), backend, tx)
	require.NoError(t, err)

	// Ensure "bye" was successfully set as the greeting.
	greetingResult, err := universe.greeter.GetGreeting(nil)
	require.NoError(t, err)
	require.Equal(t, "bye", greetingResult)

	// Assert smart contract account is created and nonce incremented.
	sca, err := sca_wrapper.NewSCA(toDeployAddress, backend)
	require.NoError(t, err)
	onChainNonce, err := sca.SNonce(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1), onChainNonce)
}

func Test4337WithLinkTokenPaymaster(t *testing.T) {
	// Deploy universe.
	universe := deployTransmissionUniverse(t)
	holder1 := universe.holder1
	holder2 := universe.holder2
	backend := universe.backend

	// Ensure no greeting is already set.
	initialGreeting, err := universe.greeter.GetGreeting(nil)
	require.NoError(t, err)
	require.Equal(t, "", initialGreeting)

	// Get the address at which the Smart Contract Account will be deployed.
	toDeployAddress, err := universe.helper.CalculateSmartContractAccountAddress(
		nil,
		holder1.From,
		universe.entryPointAddress,
		universe.factoryAddress,
	)
	require.NoError(t, err)
	t.Log("Smart Contract Account Address:", toDeployAddress)

	// Get the initialization code for the Smart Contract Account.
	fullInitializeCode, err := universe.helper.GetInitCode(nil, universe.factoryAddress, holder1.From, universe.entryPointAddress)
	require.NoError(t, err)
	t.Log("Full initialization code:", common.Bytes2Hex(fullInitializeCode))

	// Construct calldata for setGreeting.
	encodedGreetingCall, err := greeterABI.Pack("setGreeting", "bye")
	require.NoError(t, err)
	t.Log("Encoded greeting call:", common.Bytes2Hex(encodedGreetingCall))

	// Construct the calldata to be passed in the user operation.
	var (
		value    = big.NewInt(0)
		nonce    = big.NewInt(0)
		deadline = big.NewInt(1000)
	)
	fullEncoding, err := universe.helper.GetFullEndTxEncoding(nil, universe.greeterAddress, value, deadline, encodedGreetingCall)
	require.NoError(t, err)
	t.Log("Full user operation calldata:", common.Bytes2Hex(fullEncoding))

	// Deposit to LINK paymaster.
	linkTokenAddress, _, linkToken, err := link_token_interface.DeployLinkToken(holder1, backend)
	require.NoError(t, err)
	linkEthFeedAddress, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
		holder1,
		backend,
		18,
		(*big.Int)(assets.GWei(5000000)), // .005 ETH
	)
	require.NoError(t, err)
	paymasterAddress, _, _, err := paymaster_wrapper.DeployPaymaster(holder1, backend, linkTokenAddress, linkEthFeedAddress, universe.entryPointAddress)
	require.NoError(t, err)
	backend.Commit()
	tx, err := linkToken.TransferAndCall(
		holder1,
		paymasterAddress,
		assets.Ether(1000).ToInt(),
		common.LeftPadBytes(toDeployAddress.Bytes(), 32),
	)
	require.NoError(t, err)
	backend.Commit()
	_, err = bind.WaitMined(testutils.Context(t), backend, tx)
	require.NoError(t, err)

	// Construct and execute user operation.
	userOp := entry_point.UserOperation{
		Sender:               toDeployAddress,
		Nonce:                nonce,
		InitCode:             fullInitializeCode,
		CallData:             fullEncoding,
		CallGasLimit:         big.NewInt(10_000_000),
		VerificationGasLimit: big.NewInt(10_000_000),
		PreVerificationGas:   big.NewInt(10_000_000),
		MaxFeePerGas:         big.NewInt(100),
		MaxPriorityFeePerGas: big.NewInt(200),
		PaymasterAndData:     paymasterAddress.Bytes(),
		Signature:            []byte(""),
	}

	// Generate hash from user operation, sign it, and include it in the user operation.
	userOpHash, err := universe.entryPoint.GetUserOpHash(nil, userOp)
	require.NoError(t, err)
	fullHash, err := universe.helper.GetFullHashForSigning(nil, userOpHash, toDeployAddress)
	require.NoError(t, err)
	t.Log("Full hash for signing:", common.Bytes2Hex(fullHash[:]))
	sig, err := transmission.SignMessage(universe.holder1Key.ToEcdsaPrivKey(), fullHash[:])
	require.NoError(t, err)
	t.Log("Signature:", common.Bytes2Hex(sig))
	userOp.Signature = sig

	// Deposit to the Paymaster's account to pay for this transaction.
	holder1.Value = assets.Ether(10).ToInt()
	tx, err = universe.entryPoint.DepositTo(holder1, paymasterAddress)
	require.NoError(t, err)
	backend.Commit()
	_, err = bind.WaitMined(testutils.Context(t), backend, tx)
	require.NoError(t, err)
	holder1.Value = assets.Ether(0).ToInt()
	balance, err := universe.entryPoint.BalanceOf(nil, paymasterAddress)
	require.NoError(t, err)
	require.Equal(t, assets.Ether(10).ToInt(), balance)

	// Run handleOps from holder2's account, to demonstrate that any account can execute this signed user operation.
	tx, err = universe.entryPoint.HandleOps(holder2, []entry_point.UserOperation{userOp}, holder1.From)
	require.NoError(t, err)
	backend.Commit()
	_, err = bind.WaitMined(testutils.Context(t), backend, tx)
	require.NoError(t, err)

	// Ensure "bye" was successfully set as the greeting.
	greetingResult, err := universe.greeter.GetGreeting(nil)
	require.NoError(t, err)
	require.Equal(t, "bye", greetingResult)

	// Assert smart contract account is created and nonce incremented.
	sca, err := sca_wrapper.NewSCA(toDeployAddress, backend)
	require.NoError(t, err)
	onChainNonce, err := sca.SNonce(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1), onChainNonce)
}

func Test4337WithLinkTokenVRFRequestAndPaymaster(t *testing.T) {
	// Deploy universe.
	universe := deployTransmissionUniverse(t)
	holder1 := universe.holder1
	holder2 := universe.holder2
	backend := universe.backend

	// Get the address at which the Smart Contract Account will be deployed.
	toDeployAddress, err := universe.helper.CalculateSmartContractAccountAddress(
		nil,
		holder1.From,
		universe.entryPointAddress,
		universe.factoryAddress,
	)
	require.NoError(t, err)
	t.Log("Smart Contract Account Address:", toDeployAddress)

	// Get the initialization code for the Smart Contract Account.
	fullInitializeCode, err := universe.helper.GetInitCode(nil, universe.factoryAddress, holder1.From, universe.entryPointAddress)
	require.NoError(t, err)
	t.Log("Full initialization code:", common.Bytes2Hex(fullInitializeCode))

	// Construct calldata for the vrf request.
	var keyhash [32]byte
	copy(keyhash[:], common.LeftPadBytes(big.NewInt(123).Bytes(), 32))
	var fee = assets.Ether(1).ToInt()
	encodedVRFRequest, err := consumerABI.Pack("doRequestRandomness", keyhash, fee)
	require.NoError(t, err)
	t.Log("Encoded vrf request:", common.Bytes2Hex(encodedVRFRequest))

	// Construct the calldata to be passed in the user operation.
	var (
		value    = big.NewInt(0)
		nonce    = big.NewInt(0)
		deadline = big.NewInt(1000)
	)
	fullEncoding, err := universe.helper.GetFullEndTxEncoding(nil, universe.vrfConsumerAddress, value, deadline, encodedVRFRequest)
	require.NoError(t, err)
	t.Log("Full user operation calldata:", common.Bytes2Hex(fullEncoding))

	// Deposit to LINK paymaster.
	paymasterAddress, _, _, err := paymaster_wrapper.DeployPaymaster(holder1, backend, universe.linkTokenAddress, universe.linkEthFeedAddress, universe.entryPointAddress)
	require.NoError(t, err)
	backend.Commit()
	tx, err := universe.linkToken.TransferAndCall(
		holder1,
		paymasterAddress,
		assets.Ether(1000).ToInt(),
		common.LeftPadBytes(toDeployAddress.Bytes(), 32),
	)
	require.NoError(t, err)
	backend.Commit()
	_, err = bind.WaitMined(testutils.Context(t), backend, tx)
	require.NoError(t, err)

	// Generate encoded paymaster data to fund the VRF consumer.
	encodedPaymasterData, err := universe.helper.GetAbiEncodedDirectRequestData(nil, universe.vrfConsumerAddress, fee, fee)
	require.NoError(t, err)

	// Construct and execute user operation.
	userOp := entry_point.UserOperation{
		Sender:               toDeployAddress,
		Nonce:                nonce,
		InitCode:             fullInitializeCode,
		CallData:             fullEncoding,
		CallGasLimit:         big.NewInt(10_000_000),
		VerificationGasLimit: big.NewInt(10_000_000),
		PreVerificationGas:   big.NewInt(10_000_000),
		MaxFeePerGas:         big.NewInt(100),
		MaxPriorityFeePerGas: big.NewInt(200),
		PaymasterAndData:     append(append(paymasterAddress.Bytes(), byte(0)), encodedPaymasterData...),
		Signature:            []byte(""),
	}

	// Generate hash from user operation, sign it, and include it in the user operation.
	userOpHash, err := universe.entryPoint.GetUserOpHash(nil, userOp)
	require.NoError(t, err)
	fullHash, err := universe.helper.GetFullHashForSigning(nil, userOpHash, toDeployAddress)
	require.NoError(t, err)
	t.Log("Full hash for signing:", common.Bytes2Hex(fullHash[:]))
	sig, err := transmission.SignMessage(universe.holder1Key.ToEcdsaPrivKey(), fullHash[:])
	require.NoError(t, err)
	t.Log("Signature:", common.Bytes2Hex(sig))
	userOp.Signature = sig

	// Deposit to the Paymaster's account to pay for this transaction.
	holder1.Value = assets.Ether(10).ToInt()
	tx, err = universe.entryPoint.DepositTo(holder1, paymasterAddress)
	require.NoError(t, err)
	backend.Commit()
	_, err = bind.WaitMined(testutils.Context(t), backend, tx)
	require.NoError(t, err)
	holder1.Value = assets.Ether(0).ToInt()
	balance, err := universe.entryPoint.BalanceOf(nil, paymasterAddress)
	require.NoError(t, err)
	require.Equal(t, assets.Ether(10).ToInt(), balance)

	// Run handleOps from holder2's account, to demonstrate that any account can execute this signed user operation.
	// Manually execute transaction to test ABI packing.
	gasPrice, err := backend.SuggestGasPrice(testutils.Context(t))
	require.NoError(t, err)
	accountNonce, err := backend.PendingNonceAt(testutils.Context(t), holder2.From)
	require.NoError(t, err)
	payload, err := entrypointABI.Pack("handleOps", []entry_point.UserOperation{userOp}, holder1.From)
	require.NoError(t, err)
	gas, err := backend.EstimateGas(testutils.Context(t), ethereum.CallMsg{
		From:     holder2.From,
		To:       &universe.entryPointAddress,
		Gas:      0,
		Data:     payload,
		GasPrice: gasPrice,
	})
	unsigned := types.NewTx(&types.LegacyTx{
		Nonce:    accountNonce,
		Gas:      gas,
		To:       &universe.entryPointAddress,
		Value:    big.NewInt(0),
		Data:     payload,
		GasPrice: gasPrice,
	})
	require.NoError(t, err)
	signedtx, err := holder2.Signer(holder2.From, unsigned)
	require.NoError(t, err)
	err = backend.SendTransaction(testutils.Context(t), signedtx)
	require.NoError(t, err)
	backend.Commit()
	receipt, err := bind.WaitMined(testutils.Context(t), backend, signedtx)
	require.NoError(t, err)
	t.Log("Receipt:", receipt.Status)

	// Assert the VRF request was correctly made.
	logs, err := backend.FilterLogs(testutils.Context(t), ethereum.FilterQuery{
		Addresses: []common.Address{universe.vrfCoordinatorAddress},
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(logs))
	randomnessRequestLog, err := universe.vrfCoordinator.ParseRandomnessRequest(logs[0])
	require.NoError(t, err)
	require.Equal(t, fee, randomnessRequestLog.Fee)
	require.Equal(t, keyhash, randomnessRequestLog.KeyHash)
	require.Equal(t, universe.vrfConsumerAddress, randomnessRequestLog.Sender)

	// Assert smart contract account is created and nonce incremented.
	sca, err := sca_wrapper.NewSCA(toDeployAddress, backend)
	require.NoError(t, err)
	onChainNonce, err := sca.SNonce(nil)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1), onChainNonce)
}
