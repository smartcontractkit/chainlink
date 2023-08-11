package config

import (
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type V2 interface {
	Enabled() bool
	AnnounceAddresses() []string
	DefaultBootstrappers() (locators []ocrcommontypes.BootstrapperLocator)
	DeltaDial() models.Duration
	DeltaReconcile() models.Duration
	ListenAddresses() []string
}
