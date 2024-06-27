package docs

import (
	"log"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink/cfgtest"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/chainlink/v2/core/utils/config"
)

var (
	defaults toml.Core
)

func init() {
	if err := cfgtest.DocDefaultsOnly(strings.NewReader(coreTOML), &defaults, config.DecodeTOML); err != nil {
		log.Fatalf("Failed to initialize defaults from docs: %v", err)
	}
}

func CoreDefaults() (c toml.Core) {
	c.SetFrom(&defaults)
	c.Database.Dialect = dialects.Postgres // not user visible - overridden for tests only
	c.Tracing.Attributes = make(map[string]string)
	return
}
