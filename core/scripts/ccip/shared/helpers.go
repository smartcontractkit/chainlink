package shared

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
)

const RetryTiming = 5 * time.Second
const CrossChainTimout = 5 * time.Minute
const TxInclusionTimout = 3 * time.Minute

// WaitForMined wait for a tx to be included on chain. It will panic when
// the tx is reverted/successful based on the shouldSucceed parameter.
func WaitForMined(lggr logger.Logger, client ethereum.TransactionReader, hash common.Hash, shouldSucceed bool) error {
	maxIterations := TxInclusionTimout / RetryTiming
	for i := 0; i < int(maxIterations); i++ {
		lggr.Info("[MINING] waiting for tx to be mined...")
		receipt, _ := client.TransactionReceipt(context.Background(), hash)

		if receipt != nil {
			if shouldSucceed && receipt.Status == 0 {
				lggr.Infof("[MINING] ERROR tx reverted %s", hash.Hex())
				panic(receipt)
			} else if !shouldSucceed && receipt.Status != 0 {
				lggr.Infof("[MINING] ERROR expected tx to revert %s", hash.Hex())
				panic(receipt)
			}
			lggr.Infof("[MINING] tx mined %s successful %t", hash.Hex(), shouldSucceed)
			return nil
		}

		time.Sleep(RetryTiming)
	}
	return errors.New("No tx found within the given timeout")
}

func RequireNoError(t *testing.T, err error) {
	if err != nil {
		jErr, _ := evmclient.ExtractRPCError(err)
		t.Log(jErr)
	}
	require.NoError(t, err)
}
