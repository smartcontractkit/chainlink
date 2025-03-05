package txm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-integrations/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

type StuckTxDetectorConfig struct {
	BlockTime             time.Duration
	StuckTxBlockThreshold uint32
	DetectionURL          string
	DualBroadcast         bool
}

type stuckTxDetector struct {
	lggr         logger.Logger
	chainType    chaintype.ChainType
	config       StuckTxDetectorConfig
	lastPurgeMap map[common.Address]time.Time
}

func NewStuckTxDetector(lggr logger.Logger, chaintype chaintype.ChainType, config StuckTxDetectorConfig) *stuckTxDetector {
	return &stuckTxDetector{
		lggr:         lggr,
		chainType:    chaintype,
		config:       config,
		lastPurgeMap: make(map[common.Address]time.Time),
	}
}

func (s *stuckTxDetector) DetectStuckTransaction(ctx context.Context, tx *types.Transaction) (bool, error) {
	//nolint:gocritic //placeholder for upcoming chaintypes
	switch s.chainType {
	default:
		return s.timeBasedDetection(tx), nil
	}
}

// timeBasedDetection marks a transaction if all the following conditions are met:
// - LastBroadcastAt is not nil
// - Time since last broadcast is above the threshold
// - Time since last purge is above threshold
//
// NOTE: Potentially we can use a subset of threhsold for last purge check, because the transaction would have already been broadcasted to the mempool
// so it is more likely to be picked up compared to a transaction that hasn't been broadcasted before. This would avoid slowing down TXM for sebsequent transactions
// in case the current one is stuck.
func (s *stuckTxDetector) timeBasedDetection(tx *types.Transaction) bool {
	threshold := (s.config.BlockTime * time.Duration(s.config.StuckTxBlockThreshold))
	if tx.LastBroadcastAt != nil && min(time.Since(*tx.LastBroadcastAt), time.Since(s.lastPurgeMap[tx.FromAddress])) > threshold {
		s.lggr.Debugf("TxID: %v last broadcast was: %v and last purge: %v which is more than the max configured duration: %v. Transaction is now considered stuck and will be purged.",
			tx.ID, tx.LastBroadcastAt, s.lastPurgeMap[tx.FromAddress], threshold)
		s.lastPurgeMap[tx.FromAddress] = time.Now()
		return true
	}
	return false
}

type APIResponse struct {
	Status string      `json:"status,omitempty"`
	Hash   common.Hash `json:"hash,omitempty"`
}

const (
	APIStatusPending   = "PENDING"
	APIStatusIncluded  = "INCLUDED"
	APIStatusFailed    = "FAILED"
	APIStatusCancelled = "CANCELLED"
	APIStatusUnknown   = "UNKNOWN"
)

// Deprecated: DualBroadcastDetection doesn't provide any significant benefits in terms of speed and time
// based detection can replace it.
func (s *stuckTxDetector) DualBroadcastDetection(ctx context.Context, tx *types.Transaction) (bool, error) {
	for _, attempt := range tx.Attempts {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.config.DetectionURL+attempt.Hash.String(), nil)
		if err != nil {
			return false, fmt.Errorf("failed to make request for txID: %v, attemptHash: %v - %w", tx.ID, attempt.Hash, err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return false, fmt.Errorf("failed to get transaction status for txID: %v, attemptHash: %v - %w", tx.ID, attempt.Hash, err)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return false, fmt.Errorf("request %v failed with status: %d", req, resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return false, err
		}

		var apiResponse APIResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal response for txID: %v, attemptHash: %v - %w: %s", tx.ID, attempt.Hash, err, string(body))
		}
		switch apiResponse.Status {
		case APIStatusPending, APIStatusIncluded:
			return false, nil
		case APIStatusFailed, APIStatusCancelled:
			s.lggr.Debugf("TxID: %v with attempHash: %v was marked as failed/cancelled by the RPC. Transaction is now considered stuck and will be purged.",
				tx.ID, attempt.Hash)
			return true, nil
		case APIStatusUnknown:
			continue
		default:
			continue
		}
	}
	return false, nil
}
