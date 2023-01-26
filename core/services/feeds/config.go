package feeds

import (
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --quiet --name Config --output ./mocks/ --case=underscore

type Config interface {
	pg.QConfig

	Dev() bool
	FeatureOffchainReporting() bool
	DefaultHTTPTimeout() models.Duration
}
