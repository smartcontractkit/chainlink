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

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/types"
)

type StuckTxDetectorConfig struct {
	BlockTime             time.Duration
	StuckTxBlockThreshold uint32
	DetectionURL          string
}

type stuckTxDetector struct {
	lggr      logger.Logger
	chainType chaintype.ChainType
	config    StuckTxDetectorConfig
}

func NewStuckTxDetector(lggr logger.Logger, chaintype chaintype.ChainType, config StuckTxDetectorConfig) *stuckTxDetector {
	return &stuckTxDetector{
		lggr:      lggr,
		chainType: chaintype,
		config:    config,
	}
}

func (s *stuckTxDetector) DetectStuckTransaction(ctx context.Context, tx *types.Transaction) (bool, error) {
	switch s.chainType {
	// TODO: rename
	case chaintype.ChainDualBroadcast:
		result, err := s.dualBroadcastDetection(ctx, tx)
		if result || err != nil {
			return result, err
		}
		return s.timeBasedDetection(tx), nil
	default:
		return s.timeBasedDetection(tx), nil
	}
}

func (s *stuckTxDetector) timeBasedDetection(tx *types.Transaction) bool {
	threshold := (s.config.BlockTime * time.Duration(s.config.StuckTxBlockThreshold))
	if time.Since(tx.LastBroadcastAt) > threshold && !tx.LastBroadcastAt.IsZero() {
		s.lggr.Debugf("TxID: %v last broadcast was: %v which is more than the max configured duration: %v. Transaction is now considered stuck and will be purged.",
			tx.ID, tx.LastBroadcastAt, threshold)
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

func (s *stuckTxDetector) dualBroadcastDetection(ctx context.Context, tx *types.Transaction) (bool, error) {
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
