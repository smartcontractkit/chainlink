package plugins

import (
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// OraclePlugin is the interface that every OCR2 plugin needs to implement to be able to run from the generic
// OCR2.Delegate ServicesForSpec method.
type OraclePlugin interface {
	// GetPluginFactory return the ocr2types.ReportingPluginFactory object for the given Plugin.
	GetPluginFactory() (plugin ocr2types.ReportingPluginFactory, err error)
	// GetServices returns any additional services that the plugin might need. This can return an empty slice when
	// there are no additional services needed.
	GetServices() (services []job.ServiceCtx, err error)
}
