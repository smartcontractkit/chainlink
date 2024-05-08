package inprocessprovider

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// RegisterStandAloneProvider register the servers needed for a plugin provider,
// this is a workaround to test the Node API POCs on EVM until the EVM relayer is loopifyed
func RegisterStandAloneProvider(s *grpc.Server, p types.PluginProvider, pType types.OCR2PluginType) error {
	switch pType {
	case types.Median:
		provider, ok := p.(types.MedianProvider)
		if !ok {
			return fmt.Errorf("expected median provider got %T", p)
		}
		relayer.RegisterStandAloneMedianProvider(s, provider)
		return nil
	case types.GenericPlugin:
		relayer.RegisterStandAlonePluginProvider(s, p)
		return nil
	case types.OCR3Capability:
		provider, ok := p.(types.OCR3CapabilityProvider)
		if !ok {
			return fmt.Errorf("expected OCR3 capability provider got %T", p)
		}
		relayer.RegisterStandAloneOCR3CapabilityProvider(s, provider)
		return nil
	default:
		return fmt.Errorf("unsupported stand alone provider: %q", pType)
	}
}
