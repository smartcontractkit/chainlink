package metatx_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/transmission/generated/entry_point"
	greeter_contract "github.com/smartcontractkit/chainlink/core/gethwrappers/transmission/generated/greeter"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/transmission/generated/smart_contract_account_factory"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/transmission/generated/smart_contract_account_helper"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/metatx"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestSCA(t *testing.T) {
	// Create a private key for holder1 that we can use to sign
	accountKey := os.Getenv("ACCOUNT_KEY")
	require.NotEmpty(t, accountKey)
	b, err := hex.DecodeString(accountKey)
	require.NoError(t, err)
	d := new(big.Int).SetBytes(b)
	pkX, pkY := crypto.S256().ScalarBaseMult(d.Bytes())
	privateKey := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: crypto.S256(),
			X:     pkX,
			Y:     pkY,
		},
		D: d,
	}
	holder1Key := ethkey.FromPrivateKey(&privateKey)
	t.Log("Holder key:", holder1Key.String())

	// Construct simulated blockchain environmnet.
	holder1Transactor, err := bind.NewKeyedTransactorWithChainID(holder1Key.ToEcdsaPrivKey(), testutils.SimulatedChainID)
	var (
		metaERC20Owner = testutils.MustNewSimTransactor(t)
		holder1        = holder1Transactor
		holder2        = testutils.MustNewSimTransactor(t)
		relay          = testutils.MustNewSimTransactor(t)
	)
	genesisData := core.GenesisAlloc{
		metaERC20Owner.From: {Balance: assets.Ether(1000).ToInt()},
		holder1.From:        {Balance: assets.Ether(1000).ToInt()},
		holder2.From:        {Balance: assets.Ether(1000).ToInt()},
		relay.From:          {Balance: assets.Ether(1000).ToInt()},
	}
	gasLimit := uint32(30e6)
	backend := cltest.NewSimulatedBackend(t, genesisData, gasLimit)
	backend.Commit()

	// Deploy Entry Point, and Smart Contract Account factory and helper.
	entryPointAddress, _, entryPoint, err := entry_point.DeployEntryPoint(holder1, backend)
	require.NoError(t, err)
	factoryAddress, _, _, _ := smart_contract_account_factory.DeploySmartContractAccountFactory(holder1, backend)
	require.NoError(t, err)
	_, _, helper, err := smart_contract_account_helper.DeploySmartContractAccountHelper(holder1, backend)
	require.NoError(t, err)
	greeterAddress, _, greeter, err := greeter_contract.DeployGreeter(holder1, backend)
	require.NoError(t, err)
	backend.Commit()

	// Ensure no greeting is already set.
	initialGreeting, err := greeter.GetGreeting(nil)
	require.NoError(t, err)
	require.Equal(t, "", initialGreeting)

	// Get the address at which the Smart Contract Account will be deployed.
	toDeployAddress, err := helper.CalculateSmartContractAccountAddress(
		nil,
		holder1.From,
		entryPointAddress,
		factoryAddress,
	)
	require.NoError(t, err)
	t.Log("Smart Contrac Account Address:", holder1Key.String())

	// Get the initialization code for the Smart Contract Account.
	fullInitializeCode, err := helper.GetInitCode(nil, factoryAddress, holder1.From, entryPointAddress)
	require.NoError(t, err)
	t.Log("Full initialization code:", common.Bytes2Hex(fullInitializeCode))

	// Construct calldata for setGreeting.
	abi, err := greeter_contract.GreeterMetaData.GetAbi()
	require.NoError(t, err)
	greeting, err := utils.ABIEncode(`[{"type":"string"}]`, "bye")
	require.NoError(t, err)
	encodedGreetingCall := append(abi.Methods["setGreeting"].ID, greeting...)
	t.Log("Encoded greeting call:", common.Bytes2Hex(encodedGreetingCall))

	// Construct the calldata to be passed in the user operation.
	var (
		value    = big.NewInt(0)
		nonce    = big.NewInt(0)
		deadline = big.NewInt(1000)
	)
	fullEncoding, err := helper.GetFullEndTxEncoding(nil, greeterAddress, value, deadline, encodedGreetingCall)
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
	userOpHash, err := entryPoint.GetUserOpHash(nil, userOp)
	require.NoError(t, err)
	fullHash, err := helper.GetFullHashForSigning(nil, userOpHash)
	require.NoError(t, err)
	t.Log("Full hash for signing:", common.Bytes2Hex(fullHash[:]))
	sig, err := metatx.SignMessage(holder1Key.ToEcdsaPrivKey(), fullHash[:])
	require.NoError(t, err)
	t.Log("Signature:", common.Bytes2Hex(sig))
	userOp.Signature = sig

	// Deposit to the SCA's account to pay for this transaction.
	holder1.Value = assets.Ether(10).ToInt()
	tx, err := entryPoint.DepositTo(holder1, toDeployAddress)
	require.NoError(t, err)
	backend.Commit()
	bind.WaitMined(nil, backend, tx)
	holder1.Value = assets.Ether(0).ToInt()
	balance, err := entryPoint.BalanceOf(nil, toDeployAddress)
	require.NoError(t, err)
	require.Equal(t, assets.Ether(10).ToInt(), balance)

	// Run handleOps from holder2's account, to demonstrate that any account can execute this signed user operation.
	tx, err = entryPoint.HandleOps(holder2, []entry_point.UserOperation{userOp}, holder1.From)
	require.NoError(t, err)
	backend.Commit()
	bind.WaitMined(nil, backend, tx)

	// Ensure "bye" was successfully set as the greeting.
	greetingResult, err := greeter.GetGreeting(nil)
	require.NoError(t, err)
	require.Equal(t, "bye", greetingResult)
}
