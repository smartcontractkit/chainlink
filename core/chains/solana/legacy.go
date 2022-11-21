package solana

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/sqlx"

	solanadb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type SetupConfig interface {
	SolanaNodes() string
	pg.QConfig
}

// SetupNodes is a hack/shim method to allow node operators to specify multiple nodes via ENV.
// See: https://app.shortcut.com/chainlinklabs/epic/33587/overhaul-config?cf_workflow=500000005&ct_workflow=all
func SetupNodes(db *sqlx.DB, cfg SetupConfig, lggr logger.Logger) (err error) {
	str := cfg.SolanaNodes()
	if str == "" {
		return nil
	}

	var nodes []solanadb.Node
	if err = json.Unmarshal([]byte(str), &nodes); err != nil {
		return errors.Wrapf(err, "invalid SOLANA_NODES json: %q", str)
	}
	// Sorting gives a consistent insert ordering
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

	lggr.Info("SOLANA_NODES was set; clobbering solana_nodes table")

	orm := NewORM(db, lggr, cfg)
	return orm.SetupNodes(nodes, uniqueIDs(nodes))
}

func uniqueIDs(ns []solanadb.Node) (ids []string) {
	m := map[string]struct{}{}
	for _, n := range ns {
		id := n.SolanaChainID
		if _, ok := m[id]; ok {
			continue
		}
		ids = append(ids, id)
		m[id] = struct{}{}
	}
	return
}
