package config

import (
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// P2PV2Networking is a subset of global config relevant to p2p v2 networking.
type P2PV2Networking interface {
	P2PV2AnnounceAddresses() []string
	P2PV2Bootstrappers() (locators []ocrcommontypes.BootstrapperLocator)
	P2PV2BootstrappersRaw() []string
	P2PV2DeltaDial() models.Duration
	P2PV2DeltaReconcile() models.Duration
	P2PV2ListenAddresses() []string
}
