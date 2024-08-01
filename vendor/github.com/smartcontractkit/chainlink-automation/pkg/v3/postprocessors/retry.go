package postprocessors

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
	"github.com/smartcontractkit/chainlink-automation/pkg/v3/types"
	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

func NewRetryablePostProcessor(q types.RetryQueue, logger *log.Logger) *retryablePostProcessor {
	return &retryablePostProcessor{
		logger: log.New(logger.Writer(), fmt.Sprintf("[%s | retryable-post-processor]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
		q:      q,
	}
}

type retryablePostProcessor struct {
	logger *log.Logger
	q      types.RetryQueue
}

var _ PostProcessor = (*retryablePostProcessor)(nil)

func (p *retryablePostProcessor) PostProcess(_ context.Context, results []ocr2keepers.CheckResult, payloads []ocr2keepers.UpkeepPayload) error {
	var err error
	retryable := 0
	for i, res := range results {
		if res.PipelineExecutionState != 0 && res.Retryable {
			e := p.q.Enqueue(types.RetryRecord{
				Payload:  payloads[i],
				Interval: res.RetryInterval,
			})
			if e == nil {
				retryable++
			}
			err = errors.Join(err, e)
		}
	}
	p.logger.Printf("post-processing %d results, %d retryable\n", len(results), retryable)
	return err
}
