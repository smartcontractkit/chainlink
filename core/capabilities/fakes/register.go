package fakes

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

const EnableFakeStreamsTriggerEnvVar = "CL_ENABLE_FAKE_STREAMS_TRIGGER"

func RegisterFakeStreamsTrigger(ctx context.Context, lggr logger.Logger, registry core.CapabilitiesRegistry, nSigners int) (*fakeStreamsTrigger, error) {
	trigger := NewFakeStreamsTrigger(lggr, nSigners)
	if err := registry.Add(ctx, trigger); err != nil {
		return nil, fmt.Errorf("add fake streams trigger: %w", err)
	}

	return trigger, nil
}
