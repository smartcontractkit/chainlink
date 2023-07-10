package docs

import (
	"log"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink/cfgtest"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
	"github.com/smartcontractkit/chainlink/v2/core/utils/config"
)

var (
	defaults   toml.Core
	defaultsMu sync.Mutex
)

func CoreDefaults() (c toml.Core) {
	defaultsMu.Lock()
	defer defaultsMu.Unlock()

	if (defaults == toml.Core{}) {
		if err := cfgtest.DocDefaultsOnly(strings.NewReader(coreDefaultsTOML), &defaults, config.DecodeTOML); err != nil {
			log.Fatalf("Failed to initialize defaults from docs: %v", err)
		}
	}

	c.SetFrom(&defaults)
	c.Database.Dialect = dialects.Postgres // not user visible - overridden for tests only
	return
}
