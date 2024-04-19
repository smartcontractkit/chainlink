package types

import (
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type ReportingPluginFactory interface {
	Service
	libocr.ReportingPluginFactory
}
