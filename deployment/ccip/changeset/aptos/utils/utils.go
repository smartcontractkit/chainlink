package utils

import (
	"fmt"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink/deployment"
)

// ConfirmTx confirms aptos transactions
//
// Optional arguments:
//   - aptos.PollPeriod: time.Duration, how often to poll for the transaction. Default 100ms.
//   - aptos.PollTimeout: time.Duration, how long to wait for the transaction. Default 10s.
func ConfirmTx(chain deployment.AptosChain, txHash string, opts ...any) error {
	userTx, err := chain.Client.WaitForTransaction(txHash, opts...)
	if err != nil {
		return err
	}
	if !userTx.Success {
		return fmt.Errorf("transaction failed: %s", userTx.VmStatus)
	}
	return nil
}

func IsMCMSStagingAreaClean(client aptos.AptosRpcClient, aptosMCMSObjAddr aptos.AccountAddress) (bool, error) {
	resources, err := client.AccountResources(aptosMCMSObjAddr)
	if err != nil {
		return false, err
	}
	for _, resource := range resources {
		if strings.Contains(resource.Type, "StagingArea") {
			return false, nil
		}
	}
	return true, nil
}
