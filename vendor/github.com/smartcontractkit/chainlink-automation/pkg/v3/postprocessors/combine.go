package postprocessors

import (
	"context"
	"errors"

	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"
)

// CombinedPostprocessor ...
type CombinedPostprocessor struct {
	source []PostProcessor
}

// NewCombinedPostprocessor ...
func NewCombinedPostprocessor(src ...PostProcessor) *CombinedPostprocessor {
	return &CombinedPostprocessor{source: src}
}

// PostProcess implements the PostProcessor interface and runs all source
// processors in the sequence in which they were provided. All processors are
// run and errors are joined.
func (cpp *CombinedPostprocessor) PostProcess(ctx context.Context, results []ocr2keepers.CheckResult, payloads []ocr2keepers.UpkeepPayload) error {
	var err error

	for _, pp := range cpp.source {
		err = errors.Join(err, pp.PostProcess(ctx, results, payloads))
	}

	return err
}
