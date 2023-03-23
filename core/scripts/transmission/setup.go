package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink/core/assets"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/solidity_vrf_consumer_interface_v08"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/transmission/generated/entry_point"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/transmission/generated/greeter_wrapper"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/transmission/generated/paymaster_wrapper"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/transmission/generated/sca_wrapper"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/transmission/generated/smart_contract_account_factory"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/transmission/generated/smart_contract_account_helper"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/services/transmission"
)

var (
	greeterABI    = evmtypes.MustGetABI(greeter_wrapper.GreeterABI)
	consumerABI   = evmtypes.MustGetABI(solidity_vrf_consumer_interface_v08.VRFConsumerABI)
	entrypointABI = evmtypes.MustGetABI(entry_point.EntryPointABI)
)

func prepareSetGreeting(e helpers.Environment) {
	cmd := flag.NewFlagSet("beacon-deploy", flag.ExitOnError)
	entryPointAddressString := cmd.String("entrypoint-address", "", "entrypoint contract address")
	smartContractAccountHelperAddressString := cmd.String("helper-address", "", "smart contract account helper contract address")
	smartContractAccountFactoryAddressString := cmd.String("factory-address", "", "smart contract account factory contract address")
	greeterAddressString := cmd.String("greeter-address", "", "greeter contract address")
	greetingString := cmd.String("greeting", "hello", "the greeting to be set")
	linkTokenAddressString := cmd.String("link-address", "", "link token contract address")
	paymasterAddressString := cmd.String("paymaster-address", "", "paymaster contract address")

	topupAmountString := cmd.String("topup-amount", "0", "amount to top up paymaster subscription")
	paymasterTopupAmountString := cmd.String("paymaster-topup-amount", "0", "amount to top up paymaster's entrypoint deposit")
	deadlineString := cmd.String("deadline", "1000000", "deadline for meta-tx")
	valueString := cmd.String("value", "0", "value to be paid for meta-tx")

	callGasLimit := cmd.Int64("call-gas-limit", 1_000_000, "end-tx gas limit")
	verificationGasLimit := cmd.Int64("verification-gas-limit", 1_000_000, "gas limit for SCA deployment & verification")
	preVerificationGas := cmd.Int64("pre-verification-gas-limit", 50_000, "extra gas for entrypoint operations")
	helpers.ParseArgs(cmd, os.Args[2:], "entrypoint-address", "helper-address")

	// Assign deployed contracts.
	entryPointAddress := common.HexToAddress(*entryPointAddressString)
	entryPoint, err := entry_point.NewEntryPoint(entryPointAddress, e.Ec)
	helpers.PanicErr(err)
	helper, err := smart_contract_account_helper.NewSmartContractAccountHelper(common.HexToAddress(*smartContractAccountHelperAddressString), e.Ec)
	helpers.PanicErr(err)

	// Deploy new contracts.
	var (
		smartContractAccountFactoryAddress common.Address = common.HexToAddress(*smartContractAccountFactoryAddressString)
		greeterAddress                     common.Address = common.HexToAddress(*greeterAddressString)
		linkTokenAddress                   common.Address = common.HexToAddress(*linkTokenAddressString)
		linkToken                          *link_token_interface.LinkToken
		paymasterAddress                   common.Address = common.HexToAddress(*paymasterAddressString)
	)

	// Deploy Smart Contract Account Factory if not provided..
	if len(*smartContractAccountFactoryAddressString) == 0 {
		fmt.Println("\nDeploying smart contract account factory...")
		address, tx, _, err := smart_contract_account_factory.DeploySmartContractAccountFactory(e.Owner, e.Ec)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
		smartContractAccountFactoryAddress = address
		// smartContractAccountFactory = factory
	}

	// Deploy Greeter contract if not provided.
	if len(*greeterAddressString) == 0 {
		fmt.Println("\nDeploying greeter...")
		address, tx, _, err := greeter_wrapper.DeployGreeter(e.Owner, e.Ec)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
		greeterAddress = address
	}

	// Deploy LINK token if not provided. Otherwise, assign link token.
	if len(*linkTokenAddressString) == 0 {
		fmt.Println("\nDeploying link token...")
		address, tx, token, err := link_token_interface.DeployLinkToken(e.Owner, e.Ec)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
		linkTokenAddress = address
		linkToken = token
	} else {
		linkToken, err = link_token_interface.NewLinkToken(linkTokenAddress, e.Ec)
		helpers.PanicErr(err)
	}

	// Deploy Paymaster if not provided.
	if len(*paymasterAddressString) == 0 {
		fmt.Println("\nDeploying paymaster...")
		address, tx, _, err := paymaster_wrapper.DeployPaymaster(e.Owner, e.Ec, linkTokenAddress)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
		paymasterAddress = address
	}

	// Get the address at which the Smart Contract Account will be deployede.
	toDeployAddress, err := helper.CalculateSmartContractAccountAddress(
		nil,
		e.Owner.From,
		entryPointAddress,
		smartContractAccountFactoryAddress,
	)
	helpers.PanicErr(err)
	fmt.Println("\nSmart Contract Account address:", toDeployAddress)

	// Derive the nonce from the Smart Contract Account address.
	nonce := big.NewInt(0)
	code, err := e.Ec.CodeAt(context.Background(), toDeployAddress, nil)
	helpers.PanicErr(err)
	if len(code) > 0 {
		sca, err := sca_wrapper.NewSCA(toDeployAddress, e.Ec)
		helpers.PanicErr(err)
		onChainNonce, err := sca.SNonce(nil)
		helpers.PanicErr(err)
		nonce = onChainNonce
	}

	// Top up LINK subscription on paymaster.
	topupAmount, ok := big.NewInt(0).SetString(*topupAmountString, 10)
	if !ok {
		helpers.PanicErr(errors.New("failed to parse top-up amount to number"))
	}
	if topupAmount.Uint64() > 0 {
		fmt.Println("\nTopping up paymaster LINK subscription...")
		tx, err := linkToken.TransferAndCall(
			e.Owner,
			paymasterAddress,
			topupAmount,
			common.LeftPadBytes(toDeployAddress.Bytes(), 32),
		)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "topping up paymaster LINK subscription")
	}

	// Get the initialization code for the Smart Contract Account and the end-greeting call data.
	fullInitializeCode, err := helper.GetInitCode(nil, smartContractAccountFactoryAddress, e.Owner.From, entryPointAddress)
	helpers.PanicErr(err)
	encodedGreetingCall, err := greeterABI.Pack("setGreeting", *greetingString)
	helpers.PanicErr(err)

	value, ok := big.NewInt(0).SetString(*valueString, 10)
	if !ok {
		helpers.PanicErr(errors.New("failed to parse value to number"))
	}
	deadline, ok := big.NewInt(0).SetString(*deadlineString, 10)
	if !ok {
		helpers.PanicErr(errors.New("failed to parse deadline to number"))
	}
	fullEncoding, err := helper.GetFullEndTxEncoding(nil, greeterAddress, value, deadline, encodedGreetingCall)
	helpers.PanicErr(err)

	// For a nonce greater than zero, omit init code.
	if nonce.Uint64() > 0 {
		fullInitializeCode = nil
	}

	// Construct user operation.
	userOp := entry_point.UserOperation{
		Sender:               toDeployAddress,
		Nonce:                nonce,
		InitCode:             fullInitializeCode,
		CallData:             fullEncoding,
		CallGasLimit:         big.NewInt(*callGasLimit),
		VerificationGasLimit: big.NewInt(*verificationGasLimit),
		PreVerificationGas:   big.NewInt(*preVerificationGas),
		MaxFeePerGas:         e.Owner.GasPrice,
		MaxPriorityFeePerGas: e.Owner.GasPrice,
		PaymasterAndData:     paymasterAddress.Bytes(),
		Signature:            []byte(""),
	}

	// Generate signature on user operation.
	userOpHash, err := entryPoint.GetUserOpHash(nil, userOp)
	helpers.PanicErr(err)
	fullHash, err := helper.GetFullHashForSigning(nil, userOpHash)
	helpers.PanicErr(err)
	accountKey, set := os.LookupEnv("ACCOUNT_KEY")
	if !set {
		panic("need account key")
	}
	b, err := hex.DecodeString(accountKey)
	helpers.PanicErr(err)
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
	sig, err := transmission.SignMessage(&privateKey, fullHash[:])
	userOp.Signature = sig

	// Top up paymaster deposit.
	paymasterTopupAmount, ok := big.NewInt(0).SetString(*paymasterTopupAmountString, 10)
	if !ok {
		helpers.PanicErr(errors.New("failed to parse paymaster top-up amount to number"))
	}
	if paymasterTopupAmount.Uint64() > 0 {
		fmt.Println("\nTopping up paymaster entrypoint deposit...")
		e.Owner.Value = paymasterTopupAmount
		tx, err := entryPoint.DepositTo(e.Owner, paymasterAddress)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "topping up paymaster entrypoint deposit")
		e.Owner.Value = assets.Ether(0).ToInt()
	}
	// Execute user operation.
	/* e.Owner.GasLimit = 10_000_000
	tx, err := entryPoint.HandleOps(e.Owner, []entry_point.UserOperation{userOp}, e.Owner.From)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "executing user operation") */

	fmt.Println(
		"\nSetup complete.",
		"\nEntry point address:", entryPointAddress,
		"\nSmart contract account factory address:", smartContractAccountFactoryAddress,
		"\nSmart contract account address:", toDeployAddress,
		"\nPaymaster address:", paymasterAddress,
		"\nGreeter address:", greeterAddress,
		"\nLink token address:", linkTokenAddress,
		"\n\nPossible subsequent prepare request:",
	)
	fmt.Println(fmt.Sprintf(
		"go run . prepare-setGreeting --entrypoint-address %s --helper-address %s --factory-address %s --greeter-address %s --link-address %s --paymaster-address %s --greeting %s",
		entryPointAddress.String(),
		*smartContractAccountHelperAddressString,
		smartContractAccountFactoryAddress.String(),
		greeterAddress.String(),
		linkTokenAddress.String(),
		paymasterAddress.String(),
		*greetingString,
	))
	fmt.Println("\nPossible Transmission request:")
	fmt.Println(fmt.Sprintf(
		template,
		userOp.Sender.String(),
		userOp.Nonce.Uint64(),
		common.Bytes2Hex(userOp.InitCode),
		common.Bytes2Hex(userOp.CallData),
		userOp.CallGasLimit.Uint64(),
		userOp.VerificationGasLimit.Uint64(),
		userOp.PreVerificationGas.Uint64(),
		userOp.MaxFeePerGas.Uint64(),
		userOp.MaxPriorityFeePerGas.Uint64(),
		common.Bytes2Hex(userOp.PaymasterAndData),
		common.Bytes2Hex(userOp.Signature),
		entryPointAddress.String(),
	))

}

func prepareVRFRequest(e helpers.Environment) {
	cmd := flag.NewFlagSet("beacon-deploy", flag.ExitOnError)
	entryPointAddressString := cmd.String("entrypoint-address", "", "entrypoint contract address")
	smartContractAccountHelperAddressString := cmd.String("helper-address", "", "smart contract account helper contract address")
	smartContractAccountFactoryAddressString := cmd.String("factory-address", "", "smart contract account factory contract address")
	linkTokenAddressString := cmd.String("link-address", "", "link token contract address")
	paymasterAddressString := cmd.String("paymaster-address", "", "paymaster contract address")

	vrfConsumerAddressString := cmd.String("consumer-address", "", "vrf consumer contract address")
	vrfCoordinatorAddressString := cmd.String("coordinator-address", "", "vrf coordinator contract address")
	keyhashString := cmd.String("key-hash", "0x0476f9a745b61ea5c0ab224d3a6e4c99f0b02fce4da01143a4f70aa80ae76e8a", "vrf key hash")
	feeString := cmd.String("fee", "1000000000000000000", "vrf fee - 1 LINK default")

	topupAmountString := cmd.String("topup-amount", "0", "amount to top up paymaster subscription")
	paymasterTopupAmountString := cmd.String("paymaster-topup-amount", "0", "amount to top up paymaster's entrypoint deposit")
	deadlineString := cmd.String("deadline", "1000000", "deadline for meta-tx")
	valueString := cmd.String("value", "0", "value to be paid for meta-tx")

	callGasLimit := cmd.Int64("call-gas-limit", 1_000_000, "end-tx gas limit")
	verificationGasLimit := cmd.Int64("verification-gas-limit", 1_000_000, "gas limit for SCA deployment & verification")
	preVerificationGas := cmd.Int64("pre-verification-gas-limit", 50_000, "extra gas for entrypoint operations")
	helpers.ParseArgs(cmd, os.Args[2:], "entrypoint-address", "helper-address", "coordinator-address")

	// Assign deployed contracts.
	entryPointAddress := common.HexToAddress(*entryPointAddressString)
	entryPoint, err := entry_point.NewEntryPoint(entryPointAddress, e.Ec)
	helpers.PanicErr(err)
	helper, err := smart_contract_account_helper.NewSmartContractAccountHelper(common.HexToAddress(*smartContractAccountHelperAddressString), e.Ec)
	helpers.PanicErr(err)

	// Deploy new contracts.
	var (
		smartContractAccountFactoryAddress common.Address = common.HexToAddress(*smartContractAccountFactoryAddressString)
		linkTokenAddress                   common.Address = common.HexToAddress(*linkTokenAddressString)
		linkToken                          *link_token_interface.LinkToken
		paymasterAddress                   common.Address = common.HexToAddress(*paymasterAddressString)
		vrfConsumerAddress                 common.Address = common.HexToAddress(*vrfConsumerAddressString)
	)

	// Deploy Smart Contract Account Factory if not provided..
	if len(*smartContractAccountFactoryAddressString) == 0 {
		fmt.Println("\nDeploying smart contract account factory...")
		address, tx, _, err := smart_contract_account_factory.DeploySmartContractAccountFactory(e.Owner, e.Ec)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
		smartContractAccountFactoryAddress = address
		// smartContractAccountFactory = factory
	}

	// Deploy LINK token if not provided. Otherwise, assign link token.
	if len(*linkTokenAddressString) == 0 {
		fmt.Println("\nDeploying link token...")
		address, tx, token, err := link_token_interface.DeployLinkToken(e.Owner, e.Ec)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
		linkTokenAddress = address
		linkToken = token
	} else {
		linkToken, err = link_token_interface.NewLinkToken(linkTokenAddress, e.Ec)
		helpers.PanicErr(err)
	}

	// Deploy Paymaster if not provided.
	if len(*paymasterAddressString) == 0 {
		fmt.Println("\nDeploying paymaster...")
		address, tx, _, err := paymaster_wrapper.DeployPaymaster(e.Owner, e.Ec, linkTokenAddress)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
		paymasterAddress = address
	}

	// Deploy VRF consumer if not provided.
	if len(*vrfConsumerAddressString) == 0 {
		fmt.Println("\nDeploying vrf consumer...")
		address, tx, _, err := solidity_vrf_consumer_interface_v08.DeployVRFConsumer(e.Owner, e.Ec, common.HexToAddress(*vrfCoordinatorAddressString), linkTokenAddress)
		helpers.PanicErr(err)
		helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
		vrfConsumerAddress = address
	}

	// Get the address at which the Smart Contract Account will be deployede.
	toDeployAddress, err := helper.CalculateSmartContractAccountAddress(
		nil,
		e.Owner.From,
		entryPointAddress,
		smartContractAccountFactoryAddress,
	)
	helpers.PanicErr(err)
	fmt.Println("\nSmart Contract Account address:", toDeployAddress)

	// Derive the nonce from the Smart Contract Account address.
	nonce := big.NewInt(0)
	code, err := e.Ec.CodeAt(context.Background(), toDeployAddress, nil)
	helpers.PanicErr(err)
	if len(code) > 0 {
		sca, err := sca_wrapper.NewSCA(toDeployAddress, e.Ec)
		helpers.PanicErr(err)
		onChainNonce, err := sca.SNonce(nil)
		helpers.PanicErr(err)
		nonce = onChainNonce
	}

	// Top up LINK subscription on paymaster.
	topupAmount, ok := big.NewInt(0).SetString(*topupAmountString, 10)
	if !ok {
		helpers.PanicErr(errors.New("failed to parse top-up amount to number"))
	}
	if topupAmount.Uint64() > 0 {
		fmt.Println("\nTopping up paymaster LINK subscription...")
		tx, err := linkToken.TransferAndCall(
			e.Owner,
			paymasterAddress,
			topupAmount,
			common.LeftPadBytes(toDeployAddress.Bytes(), 32),
		)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "topping up paymaster LINK subscription")
	}

	// Generate vrf request calldata.
	fee, ok := big.NewInt(0).SetString(*feeString, 10)
	if !ok {
		helpers.PanicErr(errors.New("failed to parse vrf fee to number"))
	}
	encodedVRFRequest, err := consumerABI.Pack("doRequestRandomness", common.HexToHash(*keyhashString), fee)
	helpers.PanicErr(err)

	// Generate encoded paymaster data to fund the VRF consumer.
	encodedPaymasterData, err := helper.GetAbiEncodedDirectRequestData(nil, vrfConsumerAddress, fee, big.NewInt(0).Mul(fee, big.NewInt(2)))
	helpers.PanicErr(err)

	// Get the initialization code for the Smart Contract Account.
	fullInitializeCode, err := helper.GetInitCode(nil, smartContractAccountFactoryAddress, e.Owner.From, entryPointAddress)
	helpers.PanicErr(err)

	value, ok := big.NewInt(0).SetString(*valueString, 10)
	if !ok {
		helpers.PanicErr(errors.New("failed to parse value to number"))
	}
	deadline, ok := big.NewInt(0).SetString(*deadlineString, 10)
	if !ok {
		helpers.PanicErr(errors.New("failed to parse deadline to number"))
	}
	fullEncoding, err := helper.GetFullEndTxEncoding(nil, vrfConsumerAddress, value, deadline, encodedVRFRequest)
	helpers.PanicErr(err)

	// For a nonce greater than zero, omit init code.
	if nonce.Uint64() > 0 {
		fullInitializeCode = nil
	}

	// Construct user operation.
	userOp := entry_point.UserOperation{
		Sender:               toDeployAddress,
		Nonce:                nonce,
		InitCode:             fullInitializeCode,
		CallData:             fullEncoding,
		CallGasLimit:         big.NewInt(*callGasLimit),
		VerificationGasLimit: big.NewInt(*verificationGasLimit),
		PreVerificationGas:   big.NewInt(*preVerificationGas),
		MaxFeePerGas:         e.Owner.GasPrice,
		MaxPriorityFeePerGas: e.Owner.GasPrice,
		PaymasterAndData:     append(append(paymasterAddress.Bytes(), byte(0)), encodedPaymasterData...),
		Signature:            []byte(""),
	}

	// Generate signature on user operation.
	userOpHash, err := entryPoint.GetUserOpHash(nil, userOp)
	helpers.PanicErr(err)
	fullHash, err := helper.GetFullHashForSigning(nil, userOpHash)
	helpers.PanicErr(err)
	accountKey, set := os.LookupEnv("ACCOUNT_KEY")
	if !set {
		panic("need account key")
	}
	b, err := hex.DecodeString(accountKey)
	helpers.PanicErr(err)
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
	sig, err := transmission.SignMessage(&privateKey, fullHash[:])
	userOp.Signature = sig

	// Top up paymaster deposit.
	paymasterTopupAmount, ok := big.NewInt(0).SetString(*paymasterTopupAmountString, 10)
	if !ok {
		helpers.PanicErr(errors.New("failed to parse paymaster top-up amount to number"))
	}
	if paymasterTopupAmount.Uint64() > 0 {
		fmt.Println("\nTopping up paymaster entrypoint deposit...")
		e.Owner.Value = paymasterTopupAmount
		tx, err := entryPoint.DepositTo(e.Owner, paymasterAddress)
		helpers.PanicErr(err)
		helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "topping up paymaster entrypoint deposit")
		e.Owner.Value = assets.Ether(0).ToInt()
	}

	// Execute user operation.
	e.Owner.GasLimit = 10_000_000
	/* tx, err := entryPoint.HandleOps(e.Owner, []entry_point.UserOperation{userOp}, e.Owner.From)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "executing user operation") */

	fmt.Println(
		"\nSetup complete.",
		"\nEntry point address:", entryPointAddress,
		"\nSmart contract account factory address:", smartContractAccountFactoryAddress,
		"\nSmart contract account address:", toDeployAddress,
		"\nPaymaster address:", paymasterAddress,
		"\nVRF Consumer address:", vrfConsumerAddress,
		"\nVRF Coordinator address:", vrfCoordinatorAddressString,
		"\nLink token address:", linkTokenAddress,
		"\n\nPossible subsequent prepare request:",
	)
	fmt.Println(fmt.Sprintf(
		"go run . prepare-vrfRequest --entrypoint-address %s --helper-address %s --factory-address %s --link-address %s --paymaster-address %s --coordinator-address %s --consumer-address %s",
		entryPointAddress.String(),
		*smartContractAccountHelperAddressString,
		smartContractAccountFactoryAddress.String(),
		linkTokenAddress.String(),
		paymasterAddress.String(),
		*vrfCoordinatorAddressString,
		vrfConsumerAddress,
	))
	fmt.Println("\nPossible Transmission request:")
	fmt.Println(fmt.Sprintf(
		template,
		userOp.Sender.String(),
		userOp.Nonce.Uint64(),
		common.Bytes2Hex(userOp.InitCode),
		common.Bytes2Hex(userOp.CallData),
		userOp.CallGasLimit.Uint64(),
		userOp.VerificationGasLimit.Uint64(),
		userOp.PreVerificationGas.Uint64(),
		userOp.MaxFeePerGas.Uint64(),
		userOp.MaxPriorityFeePerGas.Uint64(),
		common.Bytes2Hex(userOp.PaymasterAndData),
		common.Bytes2Hex(userOp.Signature),
		entryPointAddress.String(),
	))
}

var template = `curl --header "Content-Type: application/json" --request POST --data '{
	"jsonrpc": "2.0",
	"id": 10,
	"method": "eth_sendUserOperation",
	"params": [
	  {
		"sender": "%s",
		"nonce": "%d",
		"initCode": "%s",
		"callData": "%s",
		"callGasLimit": "%d",
		"verificationGasLimit": "%d",
		"preVerificationGas": "%d",
		"maxFeePerGas": "%d",
		"maxPriorityFeePerGas": "%d",
		"paymasterAndData": "%s",
		"signature": "%s"
	  },
	  "%s"
	]
  }' http://localhost:2020/userOperations
`
