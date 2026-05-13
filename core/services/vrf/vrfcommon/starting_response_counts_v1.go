package vrfcommon

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"

	"github.com/smartcontractkit/chainlink/v2/core/utils/safe"
)

// GetStartingResponseCountsV1 returns fulfilled-request counts keyed by raw
// request ID (as used by legacy VRF v1 meta encoding). It is shared by VRF v2
// startup logic that reads the same tx meta shape.
func GetStartingResponseCountsV1(ctx context.Context, chain legacyevm.Chain) (map[[32]byte]uint64, error) {
	respCounts := make(map[[32]byte]uint64)
	var latestBlockNum *big.Int
	err := retry.Do(func() error {
		var err2 error
		latestBlockNum, err2 = chain.Client().LatestBlockHeight(ctx)
		return err2
	}, retry.Attempts(10), retry.Delay(500*time.Millisecond))
	if err != nil {
		return nil, err
	}
	if latestBlockNum == nil {
		return nil, errors.New("LatestBlockHeight return nil block num")
	}
	confirmedBlockNum := latestBlockNum.Int64() - int64(chain.Config().EVM().FinalityDepth())
	var counts []RespCountEntry
	counts, err = GetRespCounts(ctx, chain.TxManager(), chain.Client().ConfiguredChainID(), confirmedBlockNum)
	if err != nil {
		return respCounts, nil
	}

	for _, c := range counts {
		req := strings.Replace(c.RequestID, `"`, ``, 2)
		b, err2 := hex.DecodeString(req[2:])
		if err2 != nil {
			continue
		}
		var reqID [32]byte
		copy(reqID[:], b)
		count, err3 := safe.IntToUint64(c.Count)
		if err3 != nil {
			continue
		}
		respCounts[reqID] = count
	}

	return respCounts, nil
}
