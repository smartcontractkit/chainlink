package config

import (
	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

type V2 interface {
	Enabled() bool
	AnnounceAddresses() []string
	DefaultBootstrappers() (locators []ocrcommontypes.BootstrapperLocator)
	DeltaDial() commonconfig.Duration
	DeltaReconcile() commonconfig.Duration
	ListenAddresses() []string
}
