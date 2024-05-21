package config

import (
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	ocrcommontypes "github.com/smartcontractkit/libocr/commontypes"
)

type V2 interface {
	Enabled() bool
	AnnounceAddresses() []string
	DefaultBootstrappers() (locators []ocrcommontypes.BootstrapperLocator)
	DeltaDial() commonconfig.Duration
	DeltaReconcile() commonconfig.Duration
	ListenAddresses() []string
}
