package rpcv01

import (
	"context"
)

// Events returns all events matching the given filter
func (provider *Provider) Events(ctx context.Context, filter EventFilter) (*EventsOutput, error) {
	var result EventsOutput
	if err := do(ctx, provider.c, "starknet_getEvents", &result, filter); err != nil {
		return nil, err
	}

	return &result, nil
}
