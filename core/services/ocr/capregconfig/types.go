package capregconfig

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

// OCRConfigService provides OCR configuration from CapabilitiesRegistry.
// It provides Tracker and Digester instances for each capability.
//
// DON ID is determined internally by examining which DON(s) the current node
// belongs to. When OnNewRegistry() is called, the service expects the current
// node to have at most one config for any given capability ID (across all DONs
// the node is a member of).
//
// The returned Tracker and Digester implementations dynamically switch between
// registry-based config and legacy contract-based config based on availability.
type OCRConfigService interface {
	services.Service
	registrysyncer.Listener

	// GetConfigTracker returns a ContractConfigTracker for the specified capability.
	//
	// The returned tracker dynamically selects between registry-based and legacy config:
	// - If registry config is available, it returns registry-based config
	// - Otherwise, it delegates to the provided legacyTracker (if non-nil)
	GetConfigTracker(
		capabilityID string,
		ocrConfigKey string,
		legacyTracker ocrtypes.ContractConfigTracker,
	) (ocrtypes.ContractConfigTracker, error)

	// GetConfigDigester returns an OffchainConfigDigester for the specified capability.
	//
	// The returned digester dynamically selects between registry-based and legacy config:
	// - If registry config is available, it computes digest using registry parameters
	// - Otherwise, it delegates to the provided legacyDigester (if non-nil)
	GetConfigDigester(
		capabilityID string,
		ocrConfigKey string,
		legacyDigester ocrtypes.OffchainConfigDigester,
	) (ocrtypes.OffchainConfigDigester, error)
}
