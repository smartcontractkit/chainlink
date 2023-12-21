package config

import (
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
)

type V2 interface {
	Enabled() bool
	AnnounceAddresses() []string
	DefaultBootstrappers() (locators []ocrcommontypes.BootstrapperLocator)
	DeltaDial() sqlutil.Duration
	DeltaReconcile() sqlutil.Duration
	ListenAddresses() []string
}
