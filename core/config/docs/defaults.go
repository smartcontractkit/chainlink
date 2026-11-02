package docs

import (
	"log"
	"strings"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/config/configdoc"
	pgcommon "github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

func CoreDefaults() (c toml.Core) {
	if err := configdoc.DefaultsOnly(strings.NewReader(coreTOML), &c, commonconfig.DecodeTOML); err != nil {
		log.Fatalf("Failed to initialize defaults from docs: %v", err)
	}
	c.Database.DriverName = pgcommon.DriverPostgres // not user visible - overridden for tests only
	c.Tracing.Attributes = make(map[string]string)
	return c
}
