package loop

import "context"

type ctxKey int

const (
	ctxKeyJobID ctxKey = iota
	ctxKeyJobName
	ctxKeyContractID
	ctxKeyFeedID
)

// ContextValues is a helper for passing values via a [context.Context].
type ContextValues struct {
	JobID   any
	JobName any

	ContractID any
	FeedID     any
}

// Args returns a slice of args to pass to [logger.Logger.With].
func (v *ContextValues) Args() []any {
	return []any{
		"jobID", v.JobID,
		"jobName", v.JobName,
		"contractID", v.ContractID,
		"feedID", v.FeedID,
	}
}

// ContextWithValues returns a context.Context with values set from v.
func (v *ContextValues) ContextWithValues(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, ctxKeyJobID, v.JobID)
	ctx = context.WithValue(ctx, ctxKeyJobName, v.JobName)
	ctx = context.WithValue(ctx, ctxKeyContractID, v.ContractID)
	ctx = context.WithValue(ctx, ctxKeyFeedID, v.FeedID)
	return ctx
}

// SetValues sets v to values from the ctx.
func (v *ContextValues) SetValues(ctx context.Context) {
	v.JobID = ctx.Value(ctxKeyJobID)
	v.JobName = ctx.Value(ctxKeyJobName)
	v.ContractID = ctx.Value(ctxKeyContractID)
	v.FeedID = ctx.Value(ctxKeyFeedID)
}
