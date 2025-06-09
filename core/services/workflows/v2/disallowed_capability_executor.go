package v2

import (
	"context"
	"errors"

	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
)

type DisallowedCapabilityExecutor struct{}

var _ host.CapabilityExecutor = DisallowedCapabilityExecutor{}

func (d DisallowedCapabilityExecutor) CallCapability(_ context.Context, _ *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	return nil, errors.New("capability calls cannot be made during this execution")
}
