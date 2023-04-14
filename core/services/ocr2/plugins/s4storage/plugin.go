package s4storage

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// TODO implement
func NewS4Services() ([]job.ServiceCtx, error) {
	// TODO create libocr2.Oracle based on S4ReportingPluginFactory{}
	var services []job.ServiceCtx
	return services, nil
}

func GetS4APIService() *S4APIService {
	return NewS4Service()
}
