package loop

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// RegisterStandAloneProvider register the servers needed for a plugin provider,
// this is a workaround to test the Node API medianpoc on EVM until the EVM relayer is loopifyed
func RegisterStandAloneProvider(s *grpc.Server, p types.PluginProvider, pType types.OCR2PluginType) error {
	switch pType {
	case types.Median:
		mp, ok := p.(types.MedianProvider)
		if !ok {
			return fmt.Errorf("expected median provider got %T", p)
		}
		internal.RegisterStandAloneMedianProvider(s, mp)
		return nil
	default:
		return fmt.Errorf("stand alone provider only supports median, got %q", pType)
	}
}
