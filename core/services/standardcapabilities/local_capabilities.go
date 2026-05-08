package standardcapabilities

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	capStreams "github.com/smartcontractkit/chainlink/v2/core/capabilities/streams"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func shouldRegisterMockStreamsTrigger(localCfg coreconfig.LocalCapabilities) bool {
	return localCfg != nil && localCfg.GetCapabilityConfig(capStreams.MockTriggerCapabilityID) != nil
}

func registerOptionalMockStreamsTrigger(lggr logger.Logger, localCfg coreconfig.LocalCapabilities, registry core.CapabilitiesRegistry) error {
	if !shouldRegisterMockStreamsTrigger(localCfg) {
		return nil
	}

	_, err := capStreams.RegisterMockTrigger(lggr, registry)
	return err
}
