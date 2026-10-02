package config

import (
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

type V2 interface {
	Enabled() bool
	AnnounceAddresses() []string
	DefaultBootstrappers() (locators []ocrcommontypes.BootstrapperLocator)
	DeltaDial() config.Duration
	DeltaReconcile() config.Duration
	ListenAddresses() []string
}
