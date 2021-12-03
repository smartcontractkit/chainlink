package loader

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
)

// GetChainByID fetches the chain by it's id.
func GetChainByID(ctx context.Context, id string) (*types.Chain, error) {
	ldr := For(ctx)

	thunk := ldr.ChainsByIDLoader.Load(ctx, dataloader.StringKey(id))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	chain, ok := result.(types.Chain)
	if !ok {
		return nil, errors.New("invalid type")
	}

	return &chain, nil
}

// GetNodesByChainID fetches the nodes for a chain.
func GetNodesByChainID(ctx context.Context, id string) ([]types.Node, error) {
	ldr := For(ctx)

	thunk := ldr.NodesByChainIDLoader.Load(ctx, dataloader.StringKey(id))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	nodes, ok := result.([]types.Node)
	if !ok {
		return nil, errors.New("invalid type")
	}

	return nodes, nil
}

// GetFeedsManagerByID fetches the feed manager by ID.
func GetFeedsManagerByID(ctx context.Context, id string) (*feeds.FeedsManager, error) {
	ldr := For(ctx)

	thunk := ldr.FeedsManagersByIDLoader.Load(ctx, dataloader.StringKey(id))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	mgr, ok := result.(feeds.FeedsManager)
	if !ok {
		return nil, errors.New("invalid type")
	}

	return &mgr, nil
}

// GetJobRunsByID fetches the job runs by their ID.
func GetJobRunsByIDs(ctx context.Context, ids []int64) ([]pipeline.Run, error) {
	ldr := For(ctx)

	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = stringutils.FromInt64(id)
	}

	thunk := ldr.JobRunsByIDLoader.LoadMany(ctx, dataloader.NewKeysFromStrings(strIDs))
	results, errs := thunk()
	if errs != nil {
		merr := multierr.Combine(errs...)

		return nil, errors.Wrap(merr, "errors fetching runs")
	}

	runs := []pipeline.Run{}
	for _, result := range results {
		if run, ok := result.(pipeline.Run); ok {
			runs = append(runs, run)
		}
	}

	return runs, nil
}

// GetJobProposalsByFeedsManagerID fetches the job proposals by feeds manager ID.
func GetJobProposalsByFeedsManagerID(ctx context.Context, id string) ([]feeds.JobProposal, error) {
	ldr := For(ctx)

	thunk := ldr.JobProposalsByManagerIDLoader.Load(ctx, dataloader.StringKey(id))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	jbRuns, ok := result.([]feeds.JobProposal)
	if !ok {
		return nil, errors.New("invalid type")
	}

	return jbRuns, nil
}

// GetJobByPipelineSpecID fetches the job by pipeline spec ID.
func GetJobByPipelineSpecID(ctx context.Context, id string) (*job.Job, error) {
	ldr := For(ctx)

	thunk := ldr.JobsByPipelineSpecIDLoader.Load(ctx, dataloader.StringKey(id))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	jb, ok := result.(job.Job)
	if !ok {
		return nil, errors.New("invalid type")
	}

	return &jb, nil
}
