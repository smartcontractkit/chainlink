package docs

import (
	"log"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"

	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink/cfgtest"
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
	c.Database.DriverName = pgcommon.DriverPostgres // not user visible - overridden for tests only
	c.Tracing.Attributes = make(map[string]string)
	return
}
