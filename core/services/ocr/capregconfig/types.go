package capregconfig

import (
	"context"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

// OCRConfigService provides OCR configuration from CapabilitiesRegistry.
// Implements registrysyncer.Listener to receive registry updates and provides
// Tracker and Digester instances for each capability/DON/ocrConfigKey combination.
// The returned Tracker and Digester implementations dynamically switch between
// registry-based config and legacy contract-based config based on availability.
type OCRConfigService interface {
	registrysyncer.Listener

	// GetConfigTracker returns a ContractConfigTracker for the specified capability,
	// DON, and OCR config key. Use pb.OCR3ConfigDefaultKey for single-instance capabilities.
	// The returned tracker dynamically selects between registry-based and legacy config:
	// - If registry config is available, it returns registry-based config
	// - Otherwise, it delegates to the provided legacyTracker (if non-nil)
	GetConfigTracker(
		capabilityID string,
		donID uint32,
		ocrConfigKey string,
		legacyTracker ocrtypes.ContractConfigTracker,
	) (ocrtypes.ContractConfigTracker, error)

	// GetConfigDigester returns an OffchainConfigDigester for the specified capability,
	// DON, and OCR config key. Use pb.OCR3ConfigDefaultKey for single-instance capabilities.
	// The returned digester dynamically selects between registry-based and legacy config:
	// - If registry config is available, it computes digest using registry parameters
	// - Otherwise, it delegates to the provided legacyDigester (if non-nil)
	GetConfigDigester(
		capabilityID string,
		donID uint32,
		ocrConfigKey string,
		legacyDigester ocrtypes.OffchainConfigDigester,
	) (ocrtypes.OffchainConfigDigester, error)

	Start(ctx context.Context) error
	Close() error
}
