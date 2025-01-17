package ccip

import (
	"fmt"
	"time"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

const (
	transmitted = iota
	committed
	executed
	LokiLoadLabel = "ccip_load_test"
	ErrLokiClient = "failed to create Loki client for monitoring"
	ErrLokiPush   = "failed to push metrics to Loki"
)

// todo: Have a different struct for commit/exec?
type LokiMetric struct {
	EventType      int       `json:"event_type"`
	Timestamp      time.Time `json:"timestamp"`
	GasUsed        uint64    `json:"gas_used"`
	SequenceNumber uint64    `json:"sequence_number"`
}

func SendMetricsToLoki(l logger.Logger, lc *wasp.LokiClient, updatedLabels map[string]string, metrics *LokiMetric) {
	if err := lc.HandleStruct(wasp.LabelsMapToModel(updatedLabels), time.Now(), metrics); err != nil {
		l.Error(ErrLokiPush)
	}
}

func setLokiLabels(src, dst uint64) (map[string]string, error) {
	srcChainId, err := chainselectors.GetChainIDFromSelector(src)
	if err != nil {
		return nil, err
	}
	dstChainId, err := chainselectors.GetChainIDFromSelector(dst)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"sourceEvmChainId":    fmt.Sprintf("%s", srcChainId),
		"destEvmChainId":      fmt.Sprintf("%s", dstChainId),
		"destinationSelector": fmt.Sprintf("%d", dst),
		"testType":            LokiLoadLabel,
	}, nil
}
