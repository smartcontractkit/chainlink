package medianpoc

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type DataSource struct {
	pipelineRunner core.PipelineRunnerService
	spec           string
	lggr           logger.Logger

	current bridges.BridgeMetaData
	mu      sync.RWMutex
}

func (d *DataSource) Observe(ctx context.Context, reportTimestamp ocrtypes.ReportTimestamp) (*big.Int, error) {
	md, err := bridges.MarshalBridgeMetaData(d.currentAnswer())
	if err != nil {
		d.lggr.Warnw("unable to attach metadata for run", "err", err)
	}

	// NOTE: job metadata is automatically attached by the pipeline runner service
	vars := core.Vars{
		Vars: map[string]interface{}{
			"jobRun": md,
		},
	}

	results, err := d.pipelineRunner.ExecuteRun(ctx, d.spec, vars, core.Options{})
	if err != nil {
		return nil, err
	}

	finalResults := results.FinalResults()
	if len(finalResults) == 0 {
		return nil, errors.New("pipeline execution failed: not enough results")
	}

	finalResult := finalResults[0]
	if finalResult.Error != nil {
		return nil, fmt.Errorf("pipeline execution failed: %w", finalResult.Error)
	}

	asDecimal, err := utils.ToDecimal(finalResult.Value.Val)
	if err != nil {
		return nil, errors.New("cannot convert observation to decimal")
	}

	resultAsBigInt := asDecimal.BigInt()
	d.updateAnswer(resultAsBigInt)
	return resultAsBigInt, nil
}

func (d *DataSource) currentAnswer() (*big.Int, *big.Int) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.current.LatestAnswer, d.current.UpdatedAt
}

func (d *DataSource) updateAnswer(latestAnswer *big.Int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.current = bridges.BridgeMetaData{
		LatestAnswer: latestAnswer,
		UpdatedAt:    big.NewInt(time.Now().Unix()),
	}
}

type ZeroDataSource struct{}

func (d *ZeroDataSource) Observe(ctx context.Context, reportTimestamp ocrtypes.ReportTimestamp) (*big.Int, error) {
	return new(big.Int), nil
}
