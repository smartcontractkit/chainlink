package plugins

import (
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/services/job"
)

type OraclePlugin interface {
	GetPluginFactory() (plugin ocr2types.ReportingPluginFactory, err error)
	GetServices() (services []job.Service, err error)
}
