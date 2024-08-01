package loop

import "context"

type ctxKey int

const (
	ctxKeyJobID ctxKey = iota
	ctxKeyJobName
	ctxKeyContractID
	ctxKeyFeedID
	ctxKeyTransmitterID
)

// ContextValues is a helper for passing values via a [context.Context].
type ContextValues struct {
	JobID   any
	JobName any

	ContractID    any
	FeedID        any
	TransmitterID any
}

// Args returns a slice of args to pass to [logger.Logger.With].
func (v *ContextValues) Args() (a []any) {
	if v.JobID != nil {
		a = append(a, "jobID", v.JobID)
	}
	if v.JobName != nil {
		a = append(a, "jobName", v.JobName)
	}
	if v.ContractID != nil {
		a = append(a, "contractID", v.ContractID)
	}
	if v.FeedID != nil {
		a = append(a, "feedID", v.FeedID)
	}
	if v.TransmitterID != nil {
		a = append(a, "transmitterID", v.TransmitterID)
	}
	return
}

// ContextWithValues returns a context.Context with values set from v.
func (v *ContextValues) ContextWithValues(ctx context.Context) context.Context {
	if v.JobID != nil {
		ctx = context.WithValue(ctx, ctxKeyJobID, v.JobID)
	}
	if v.JobName != nil {
		ctx = context.WithValue(ctx, ctxKeyJobName, v.JobName)
	}
	if v.ContractID != nil {
		ctx = context.WithValue(ctx, ctxKeyContractID, v.ContractID)
	}
	if v.FeedID != nil {
		ctx = context.WithValue(ctx, ctxKeyFeedID, v.FeedID)
	}
	if v.TransmitterID != nil {
		ctx = context.WithValue(ctx, ctxKeyTransmitterID, v.TransmitterID)
	}
	return ctx
}

// SetValues sets v to values from the ctx.
func (v *ContextValues) SetValues(ctx context.Context) {
	v.JobID = ctx.Value(ctxKeyJobID)
	v.JobName = ctx.Value(ctxKeyJobName)
	v.ContractID = ctx.Value(ctxKeyContractID)
	v.FeedID = ctx.Value(ctxKeyFeedID)
	v.TransmitterID = ctx.Value(ctxKeyTransmitterID)
}
