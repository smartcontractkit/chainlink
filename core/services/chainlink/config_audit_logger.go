package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type auditLoggerConfig struct {
	C audit.AuditLoggerConfig
}

func (a auditLoggerConfig) Enabled() bool {
	return *a.C.Enabled
}

func (a auditLoggerConfig) ForwardToUrl() (models.URL, error) {
	return *a.C.ForwardToUrl, nil
}

func (a auditLoggerConfig) Environment() string {
	if !build.IsProd() {
		return "develop"
	}
	return "production"
}

func (a auditLoggerConfig) JsonWrapperKey() string {
	return *a.C.JsonWrapperKey
}

func (a auditLoggerConfig) Headers() (audit.ServiceHeaders, error) {
	return *a.C.Headers, nil
}
