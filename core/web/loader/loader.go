package loader

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

type loadersKey struct{}

type Dataloader struct {
	app chainlink.Application

	ChainsByIDLoader                          *dataloader.Loader
	EthTxAttemptsByEthTxIDLoader              *dataloader.Loader
	FeedsManagersByIDLoader                   *dataloader.Loader
	FeedsManagerChainConfigsByManagerIDLoader *dataloader.Loader
	JobProposalsByManagerIDLoader             *dataloader.Loader
	JobProposalSpecsByJobProposalID           *dataloader.Loader
	JobRunsByIDLoader                         *dataloader.Loader
	JobsByExternalJobIDs                      *dataloader.Loader
	JobsByPipelineSpecIDLoader                *dataloader.Loader
	NodesByChainIDLoader                      *dataloader.Loader
	SpecErrorsByJobIDLoader                   *dataloader.Loader
}

func New(app chainlink.Application) *Dataloader {
	var (
		nodes    = &nodeBatcher{app: app}
		chains   = &chainBatcher{app: app}
		mgrs     = &feedsBatcher{app: app}
		ccfgs    = &feedsManagerChainConfigBatcher{app: app}
		jobRuns  = &jobRunBatcher{app: app}
		jps      = &jobProposalBatcher{app: app}
		jpSpecs  = &jobProposalSpecBatcher{app: app}
		jbs      = &jobBatcher{app: app}
		attmpts  = &ethTransactionAttemptBatcher{app: app}
		specErrs = &jobSpecErrorsBatcher{app: app}
	)

	return &Dataloader{
		app: app,

		ChainsByIDLoader:                          dataloader.NewBatchedLoader(chains.loadByIDs),
		EthTxAttemptsByEthTxIDLoader:              dataloader.NewBatchedLoader(attmpts.loadByEthTransactionIDs),
		FeedsManagersByIDLoader:                   dataloader.NewBatchedLoader(mgrs.loadByIDs),
		FeedsManagerChainConfigsByManagerIDLoader: dataloader.NewBatchedLoader(ccfgs.loadByManagerIDs),
		JobProposalsByManagerIDLoader:             dataloader.NewBatchedLoader(jps.loadByManagersIDs),
		JobProposalSpecsByJobProposalID:           dataloader.NewBatchedLoader(jpSpecs.loadByJobProposalsIDs),
		JobRunsByIDLoader:                         dataloader.NewBatchedLoader(jobRuns.loadByIDs),
		JobsByExternalJobIDs:                      dataloader.NewBatchedLoader(jbs.loadByExternalJobIDs),
		JobsByPipelineSpecIDLoader:                dataloader.NewBatchedLoader(jbs.loadByPipelineSpecIDs),
		NodesByChainIDLoader:                      dataloader.NewBatchedLoader(nodes.loadByChainIDs),
		SpecErrorsByJobIDLoader:                   dataloader.NewBatchedLoader(specErrs.loadByJobIDs),
	}
}

// Middleware injects the dataloader into a gin context.
func Middleware(app chainlink.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := InjectDataloader(c.Request.Context(), app)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// InjectDataloader injects the dataloader into the context.
func InjectDataloader(ctx context.Context, app chainlink.Application) context.Context {
	return context.WithValue(ctx, loadersKey{}, New(app))
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Dataloader {
	return ctx.Value(loadersKey{}).(*Dataloader)
}
