package terra

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"

	terradb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/logger"
)

type SetupConfig interface {
	TerraNodes() string
	LogSQL() bool
}

// SetupNodes is a hack/shim method to allow node operators to specify multiple nodes via ENV.
// See: https://app.shortcut.com/chainlinklabs/epic/33587/overhaul-config?cf_workflow=500000005&ct_workflow=all
func SetupNodes(db *sqlx.DB, cfg SetupConfig, lggr logger.Logger) (err error) {
	str := cfg.TerraNodes()
	if str == "" {
		return nil
	}

	var nodes []terradb.Node
	if err = json.Unmarshal([]byte(str), &nodes); err != nil {
		return errors.Wrapf(err, "invalid TERRA_NODES json: %q", str)
	}
	// Sorting gives a consistent insert ordering
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

	lggr.Info("TERRA_NODES was set; clobbering terra_nodes table")

	orm := NewORM(db, lggr, cfg)
	return orm.SetupNodes(nodes, uniqueIDs(nodes))
}

func uniqueIDs(ns []terradb.Node) (ids []string) {
	m := map[string]struct{}{}
	for _, n := range ns {
		id := n.TerraChainID
		if _, ok := m[id]; ok {
			continue
		}
		ids = append(ids, id)
		m[id] = struct{}{}
	}
	return
}
