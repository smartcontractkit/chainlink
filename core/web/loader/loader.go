package loader

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/dataloader"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

type loadersKey struct{}

type Dataloader struct {
	app chainlink.Application

	NodesByChainIDLoader          *dataloader.Loader
	ChainsByIDLoader              *dataloader.Loader
	FeedsManagersByIDLoader       *dataloader.Loader
	JobRunsByIDLoader             *dataloader.Loader
	JobsByPipelineSpecIDLoader    *dataloader.Loader
	JobProposalsByManagerIDLoader *dataloader.Loader
}

func New(app chainlink.Application) *Dataloader {
	nodes := &nodeBatcher{app: app}
	chains := &chainBatcher{app: app}
	mgrs := &feedsBatcher{app: app}
	jobRuns := &jobRunBatcher{app: app}
	jps := &jobProposalBatcher{app: app}
	jbs := &jobBatcher{app: app}

	return &Dataloader{
		app: app,

		NodesByChainIDLoader:          dataloader.NewBatchedLoader(nodes.loadByChainIDs),
		ChainsByIDLoader:              dataloader.NewBatchedLoader(chains.loadByIDs),
		FeedsManagersByIDLoader:       dataloader.NewBatchedLoader(mgrs.loadByIDs),
		JobRunsByIDLoader:             dataloader.NewBatchedLoader(jobRuns.loadByIDs),
		JobsByPipelineSpecIDLoader:    dataloader.NewBatchedLoader(jbs.loadByPipelineSpecIDs),
		JobProposalsByManagerIDLoader: dataloader.NewBatchedLoader(jps.loadByManagersIDs),
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
