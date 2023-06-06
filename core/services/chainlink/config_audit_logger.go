package chainlink

import (
	"github.com/smartcontractkit/chainlink/v2/core/build"
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type auditLoggerConfig struct {
	c v2.AuditLogger
}

func (a auditLoggerConfig) Enabled() bool {
	return *a.c.Enabled
}

func (a auditLoggerConfig) ForwardToUrl() (models.URL, error) {
	return *a.c.ForwardToUrl, nil
}

func (a auditLoggerConfig) Environment() string {
	if !build.IsProd() {
		return "develop"
	}
	return "production"
}

func (a auditLoggerConfig) JsonWrapperKey() string {
	return *a.c.JsonWrapperKey
}

func (a auditLoggerConfig) Headers() (models.ServiceHeaders, error) {
	return *a.c.Headers, nil
}
