package config

import (
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type AuditLogger interface {
	Enabled() bool
	ForwardToUrl() (config.URL, error)
	Environment() string
	JsonWrapperKey() string
	Headers() (models.ServiceHeaders, error)
}
