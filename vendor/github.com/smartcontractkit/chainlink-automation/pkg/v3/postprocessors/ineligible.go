package postprocessors

import (
	"context"
	"errors"
	"fmt"
	"log"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	"github.com/smartcontractkit/chainlink-automation/pkg/v3/telemetry"
)

type ineligiblePostProcessor struct {
	lggr         *log.Logger
	stateUpdater ocr2keepers.UpkeepStateUpdater
}

func NewIneligiblePostProcessor(stateUpdater ocr2keepers.UpkeepStateUpdater, logger *log.Logger) *ineligiblePostProcessor {
	return &ineligiblePostProcessor{
		lggr:         log.New(logger.Writer(), fmt.Sprintf("[%s | ineligible-post-processor]", telemetry.ServiceName), telemetry.LogPkgStdFlags),
		stateUpdater: stateUpdater,
	}
}

func (p *ineligiblePostProcessor) PostProcess(ctx context.Context, results []ocr2keepers.CheckResult, _ []ocr2keepers.UpkeepPayload) error {
	var merr error
	ineligible := 0
	for _, res := range results {
		if res.PipelineExecutionState == 0 && !res.Eligible {
			err := p.stateUpdater.SetUpkeepState(ctx, res, ocr2keepers.Ineligible)
			if err != nil {
				merr = errors.Join(merr, err)
				continue
			}
			ineligible++
		}
	}
	p.lggr.Printf("post-processing %d results, %d ineligible\n", len(results), ineligible)
	return merr
}
