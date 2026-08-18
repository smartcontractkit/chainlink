package capregconfig

import (
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/capabilities/libs/ocr"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

// computeConfigDigest is the digest of an OCR config held in CRE's
// CapabilitiesRegistry contract.
//
// It moved to github.com/smartcontractkit/capabilities/libs/ocr, because every
// process running or serving one of these configurations has to compute the
// same bytes and this node is no longer the only one that does: a capability
// hosted outside it is handed configurations read from the same contract, and
// two copies of this could drift into two digests and an oracle set that will
// not talk to itself.
//
// Kept here as an alias so that nothing calling it has to change.
func computeConfigDigest(
	chainID uint64,
	registryAddress string,
	capabilityID string,
	donID uint32,
	ocrConfigKey string,
	cfg *capabilitiespb.OCR3Config,
) (ocrtypes.ConfigDigest, error) {
	return ocr.ConfigDigest(chainID, registryAddress, capabilityID, donID, ocrConfigKey, cfg)
}
