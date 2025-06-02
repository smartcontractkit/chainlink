package ccip

import (
	"context"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"testing"
	"time"
)

const (
	LokiLoadLabel = "ccipv2_load_test"
	ErrLokiPush   = "failed to push metrics to Loki"
	FinalityDepth = 200
)

type LokiMetric struct {
	TransmitTime   uint64 `json:"transmit_time"`
	SequenceNumber uint64 `json:"sequence_number"`
	CommitDuration uint64 `json:"commit_duration"`
	ExecDuration   uint64 `json:"exec_duration"`
	FailedCommit   bool   `json:"failed_commit"`
	FailedExec     bool   `json:"failed_exec"`
}

// MetricsManager is used for maintaining state of different sequence numbers
// Once we've received all expected timestamps, it pushes the metrics to Loki
type MetricManager struct {
	lggr       logger.Logger
	loki       *wasp.LokiClient
	InputChan  chan messageData
	state      map[srcDstSeqNum]metricState
	blockTimes map[uint64]uint64
	testLabel  string
	ErrorChan  chan error
}

type metricState struct {
	timestamps [3]uint64
}

type srcDstSeqNum struct {
	src    uint64
	dst    uint64
	seqNum uint64
}

type messageData struct {
	eventType int
	srcDstSeqNum
	timestamp uint64
}

func NewMetricsManager(t *testing.T, l logger.Logger, overrides *ccip.LoadConfig, blockTimes map[uint64]uint64) *MetricManager {
	// initialize loki using endpoint from user defined env vars
	loki, err := wasp.NewLokiClient(wasp.NewEnvLokiConfig())
	require.NoError(t, err)
	testLabel := "default"
	if overrides.TestLabel != nil {
		testLabel = *overrides.TestLabel
	}

	return &MetricManager{
		lggr:       l,
		loki:       loki,
		InputChan:  make(chan messageData),
		state:      make(map[srcDstSeqNum]metricState),
		blockTimes: blockTimes,
		testLabel:  testLabel,
		ErrorChan:  make(chan error),
	}
}

func (mm *MetricManager) Start(ctx context.Context) {
	defer close(mm.InputChan)
	for {
		select {
		case <-ctx.Done():
			mm.lggr.Infow("received timeout, pushing remaining state to loki")
			// any remaining data in state should be pushed to loki as incomplete
			for srcDstSeqNum, metricState := range mm.state {
				commitDuration, execDuration := uint64(0), uint64(0)
				timestamps := metricState.timestamps
				if timestamps[committed] != 0 && timestamps[transmitted] != 0 {
					commitDuration =
						timestamps[committed] - timestamps[transmitted] - mm.getChainFinalityTime(srcDstSeqNum.src)
				}
				if timestamps[executed] != 0 && timestamps[committed] != 0 {
					execDuration = timestamps[executed] - timestamps[committed]
				}

				lokiLabels, err := setLokiLabels(srcDstSeqNum.src, srcDstSeqNum.dst, mm.testLabel)
				if err != nil {
					mm.lggr.Error("error setting loki labels", "error", err)
					// don't return here, we still want to push metrics to loki
				}
				SendMetricsToLoki(mm.lggr, mm.loki, lokiLabels, &LokiMetric{
					TransmitTime:   timestamps[transmitted],
					ExecDuration:   execDuration,
					CommitDuration: commitDuration,
					SequenceNumber: srcDstSeqNum.seqNum,
					FailedCommit:   timestamps[committed] == 0,
					FailedExec:     timestamps[executed] == 0,
				})
			}
			return
		case data := <-mm.InputChan:
			if _, ok := mm.state[data.srcDstSeqNum]; !ok {
				mm.state[data.srcDstSeqNum] = metricState{
					timestamps: [3]uint64{0, 0, 0},
				}
			}

			if data.seqNum == 0 {
				// seqNum of 0 indicates an error. Push nil values to loki
				lokiLabels, err := setLokiLabels(data.src, data.dst, mm.testLabel)
				if err != nil {
					mm.lggr.Error("error setting loki labels", "error", err)
				}
				SendMetricsToLoki(mm.lggr, mm.loki, lokiLabels, &LokiMetric{
					TransmitTime:   data.timestamp,
					ExecDuration:   0,
					CommitDuration: 0,
					SequenceNumber: 0,
				})
				continue
			}
			state := mm.state[data.srcDstSeqNum]
			state.timestamps[data.eventType] = data.timestamp

			mm.state[data.srcDstSeqNum] = state
			if data.eventType == executed {
				mm.lggr.Infow("new state for received seqNum is ", "dst", data.dst, "seqNum", data.seqNum, "timestamps", state.timestamps)
			}
			lokiLabels, err := setLokiLabels(data.src, data.dst, mm.testLabel)
			if err != nil {
				mm.lggr.Error("error setting loki labels", "error", err)
			}

			// only add commit and exec durations if we have correct timestamps to calculate them
			commitDuration := uint64(0)
			if state.timestamps[committed] != 0 && state.timestamps[transmitted] != 0 {
				commitDuration =
					state.timestamps[committed] - state.timestamps[transmitted] - mm.getChainFinalityTime(data.src)
			}
			execDuration := uint64(0)
			if state.timestamps[executed] != 0 && state.timestamps[committed] != 0 {
				execDuration = state.timestamps[executed] - state.timestamps[committed]
			}

			SendMetricsToLoki(mm.lggr, mm.loki, lokiLabels, &LokiMetric{
				TransmitTime:   state.timestamps[transmitted],
				ExecDuration:   execDuration,
				CommitDuration: commitDuration,
				SequenceNumber: data.seqNum,
			})

			if state.timestamps[transmitted] != 0 && state.timestamps[committed] != 0 && state.timestamps[executed] != 0 {
				// We have a fully completed sequence number, remove it from state for subsequent tests
				delete(mm.state, data.srcDstSeqNum)
			}
		}
	}
}

func SendMetricsToLoki(l logger.Logger, lc *wasp.LokiClient, updatedLabels map[string]string, metrics *LokiMetric) {
	if err := lc.HandleStruct(wasp.LabelsMapToModel(updatedLabels), time.Now(), metrics); err != nil {
		l.Error(ErrLokiPush)
	}
}

func (mm MetricManager) getChainFinalityTime(chainSelector uint64) uint64 {
	if blockTime, ok := mm.blockTimes[chainSelector]; ok {
		return blockTime * FinalityDepth
	}
	mm.lggr.Error("block time not found for chainSelector", "chainSelector", chainSelector)
	return 0
}
func setLokiLabels(src, dst uint64, testLabel string) (map[string]string, error) {
	srcChainID, err := chainselectors.GetChainIDFromSelector(src)
	if err != nil {
		return nil, err
	}
	dstChainID, err := chainselectors.GetChainIDFromSelector(dst)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"sourceEvmChainId": srcChainID,
		"destEvmChainId":   dstChainID,
		"testType":         LokiLoadLabel,
		"testLabel":        testLabel,
	}, nil
}
