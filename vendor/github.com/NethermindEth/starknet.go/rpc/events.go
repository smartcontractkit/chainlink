package rpc

import (
	"context"
)

// Events retrieves events from the provider matching the given filter.
//
// Parameters:
// - ctx: The context to use for the request
// - input: The input parameters for retrieving events
// Returns
// - eventChunk: The retrieved events
// - error: An error if any
func (provider *Provider) Events(ctx context.Context, input EventsInput) (*EventChunk, error) {
	var result EventChunk
	if err := do(ctx, provider.c, "starknet_getEvents", &result, input); err != nil {
		return nil, tryUnwrapToRPCErr(err, ErrPageSizeTooBig, ErrInvalidContinuationToken, ErrBlockNotFound, ErrTooManyKeysInFilter)
	}
	return &result, nil
}
