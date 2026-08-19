package config

import (
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type AuditLogger interface {
	Enabled() bool
	ForwardToURL() (commonconfig.URL, error)
	Environment() string
	JSONWrapperKey() string
	Headers() (models.ServiceHeaders, error)
}
