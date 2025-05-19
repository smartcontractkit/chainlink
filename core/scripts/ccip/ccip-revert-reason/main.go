package main

import (
	"flag"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/revert-reason/handler"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/secrets"
)

var (
	errorCodeString = flag.String("errorCode", "", "Error code string (e.g. 0x08c379a0)")

	chainID     = flag.Uint64("chainId", 0, "Chain ID for the transaction (e.g. 420)")
	txHash      = flag.String("txHash", "", "Transaction hash (e.g. 0x97be8559164442595aba46b5f849c23257905b78e72ee43d9b998b28eee78b84)")
	txRequester = flag.String("txRequester", "", "Transaction requester address (e.g. 0xe88ff73814fb891bb0e149f5578796fa41f20242)")
	rpcURL      = flag.String("rpcURL", "", "RPC URL for the chain (can also be set in env var RPC_<chain_id>)")
)

func main() {
	flag.Usage = func() {
		fmt.Println("Usage: go run . [flags]")
		fmt.Println("You must provide either an error code string or the transaction details and an RPC URL to decode the error")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *errorCodeString == "" && (*chainID == 0 || *txHash == "" || *txRequester == "") {
		flag.Usage()
		return
	}

	errorString, err := getErrorString()
	if err != nil {
		fmt.Printf("Error getting error string: %v\n", err)
		return
	}
	decodedError, err := handler.DecodeErrorStringFromABI(errorString)
	if err != nil {
		fmt.Printf("Error decoding error string: %v\n", err)
		return
	}

	fmt.Println(decodedError)
}

func getErrorString() (string, error) {
	if *errorCodeString != "" {
		return *errorCodeString, nil
	}

	if *rpcURL == "" {
		fmt.Printf("RPC URL not provided, looking for RPC_%d env var\n", *chainID)
		envRPC := secrets.GetRPC(*chainID)
		rpcURL = &envRPC
	}

	ec, err := ethclient.Dial(*rpcURL)
	if err != nil {
		return "", err
	}
	errString, err := handler.GetErrorForTx(ec, *txHash, *txRequester)
	if err != nil {
		return "", err
	}

	return errString, nil
}
