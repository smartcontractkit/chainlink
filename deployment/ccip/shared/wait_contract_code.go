package shared

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sethvargo/go-retry"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
)

// WaitForContractCode retries until bytecode is visible at addr. Geth in Docker CI can
// return a mined receipt before eth_getCode serves the deployment on all RPC paths.
func WaitForContractCode(ctx context.Context, client cldf_evm.OnchainClient, addr common.Address) error {
	return retry.Do(ctx, retry.WithMaxDuration(30*time.Second, retry.WithCappedDuration(2*time.Second, retry.NewFibonacci(500*time.Millisecond))), func(ctx context.Context) error {
		code, err := client.CodeAt(ctx, addr, nil)
		if err != nil {
			return retry.RetryableError(err)
		}
		if len(code) == 0 {
			return retry.RetryableError(fmt.Errorf("no contract code at %s yet", addr))
		}
		return nil
	})
}
