package ccip

import (
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"strconv"
	"time"
)

const (
	LokiLoadLabel = "ccipv2_load_test"
	ErrLokiPush   = "failed to push metrics to Loki"
)

// MetricsManager is used for maintaining state of different sequence numbers
// Once we've received all expected timestamps, it pushes the metrics to Loki
type MetricManager struct {
	lggr      logger.Logger
	loki      *wasp.LokiClient
	InputChan chan messageData
	state     map[srcDstSeqNum][3]uint64
	DoneChan  chan struct{}
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

func NewMetricsManager(l logger.Logger, loki *wasp.LokiClient) *MetricManager {
	return &MetricManager{
		lggr:      l,
		loki:      loki,
		InputChan: make(chan messageData),
		DoneChan:  make(chan struct{}),
		state:     make(map[srcDstSeqNum][3]uint64),
	}
}

func (mm *MetricManager) Stop() {
	if isClosed(mm.DoneChan) {
		return
	}
	close(mm.DoneChan)
	close(mm.InputChan)
}

func isClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
func (mm *MetricManager) Start() {
	for {
		select {
		case <-mm.DoneChan:
			// any remaining data in state should be pushed to loki as incomplete
			for srcDstSeqNum, timestamps := range mm.state {
				lokiLabels, err := setLokiLabels(srcDstSeqNum.src, srcDstSeqNum.dst)
				if err != nil {
					mm.lggr.Error("error setting loki labels", "error", err)
					// don't return here, we still want to push metrics to loki
				}
				commitDuration, execDuration := uint64(0), uint64(0)
				if timestamps[1] != 0 && timestamps[0] != 0 {
					commitDuration = timestamps[1] - timestamps[0]
				}
				if timestamps[2] != 0 && timestamps[1] != 0 {
					execDuration = timestamps[2] - timestamps[1]
				}
				SendMetricsToLoki(mm.lggr, mm.loki, lokiLabels, &LokiMetric{
					ExecDuration:   execDuration,
					CommitDuration: commitDuration,
					SequenceNumber: srcDstSeqNum.seqNum,
				})
			}

			return
		case data := <-mm.InputChan:
			if _, ok := mm.state[data.srcDstSeqNum]; !ok {
				mm.state[data.srcDstSeqNum] = [3]uint64{0, 0, 0}
			}

			timestamps := mm.state[data.srcDstSeqNum]
			timestamps[data.eventType] = data.timestamp
			mm.state[data.srcDstSeqNum] = timestamps

			// we have all data needed to push to Loki
			if mm.state[data.srcDstSeqNum][0] != 0 && mm.state[data.srcDstSeqNum][1] != 0 && mm.state[data.srcDstSeqNum][2] != 0 {
				lokiLabels, err := setLokiLabels(data.src, data.dst)
				if err != nil {
					mm.lggr.Error("error setting loki labels", "error", err)
				}
				mm.lggr.Infow("publishing data to for ", "dst", data.dst, "src", data.src, "seqNum", data.seqNum)
				SendMetricsToLoki(mm.lggr, mm.loki, lokiLabels, &LokiMetric{
					ExecDuration:   mm.state[data.srcDstSeqNum][2] - mm.state[data.srcDstSeqNum][1],
					CommitDuration: mm.state[data.srcDstSeqNum][1] - mm.state[data.srcDstSeqNum][0],
					SequenceNumber: data.seqNum,
				})

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

func setLokiLabels(src, dst uint64) (map[string]string, error) {
	srcChainID, err := chainselectors.GetChainIDFromSelector(src)
	if err != nil {
		return nil, err
	}
	dstChainID, err := chainselectors.GetChainIDFromSelector(dst)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"sourceEvmChainId":    srcChainID,
		"sourceSelector":      strconv.FormatUint(src, 10),
		"destEvmChainId":      dstChainID,
		"destinationSelector": strconv.FormatUint(dst, 10),
		"testType":            LokiLoadLabel,
	}, nil
}
