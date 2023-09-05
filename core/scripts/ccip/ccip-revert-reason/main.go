package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/revert-reason/handler"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/secrets"
)

// How to use
// Set either an error code string OR set the chainId, txHash and txRequester.
// Setting an error code allows the script to run offline and doesn't require any RPC
// endpoint. Using the chainId, txHash and txRequester requires an RPC endpoint, and if
// the tx is old, the node needs to run in archive mode.
//
// Set the variable(s) and run main.go. The script will try to match the error code to the
// ABIs of various CCIP contracts. If it finds a match, it will check if it's a CCIP wrapped error
// like ExecutionError and TokenRateLimitError, and if so, it will decode the inner error.
//
// To configure an RPC endpoint, set the RPC_<chain_id> environment variable to the RPC endpoint.
// e.g. RPC_420=https://rpc.<chain_id>.com
const (
	ErrorCodeString = "0x4e487b710000000000000000000000000000000000000000000000000000000000000032"

	// The following inputs are only used if ERROR_CODE_STRING is empty
	// Need a node URL
	// NOTE: this node needs to run in archive mode if the tx is old
	ChainId     = uint64(420)
	TxHash      = "0x97be8559164442595aba46b5f849c23257905b78e72ee43d9b998b28eee78b84"
	TxRequester = "0xe88ff73814fb891bb0e149f5578796fa41f20242"
	EnvFileName = ".env"
)

func main() {
	errorString, err := getErrorString()
	if err != nil {
		panic(err)
	}
	decodedError, err := handler.DecodeErrorStringFromABI(errorString)
	if err != nil {
		panic(err)
	}

	fmt.Println(decodedError)
}

func getErrorString() (string, error) {
	errorCodeString := ErrorCodeString

	if errorCodeString == "" {
		// Try to load env vars from .env file
		err := godotenv.Load(EnvFileName)
		if err != nil {
			fmt.Println("No .env file found, using env vars from shell")
		}

		ec, err := ethclient.Dial(secrets.GetRPC(ChainId))
		if err != nil {
			return "", err
		}
		errorCodeString, err = handler.GetErrorForTx(ec, TxHash, TxRequester)
		if err != nil {
			return "", err
		}
	}

	return errorCodeString, nil
}
