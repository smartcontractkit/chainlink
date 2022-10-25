package starknet

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/sqlx"

	starknetdb "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"
	"github.com/smartcontractkit/chainlink/core/services/pg"

	"github.com/smartcontractkit/chainlink/core/logger"
)

type SetupConfig interface {
	StarkNetNodes() string
	pg.QConfig
}

// SetupNodes is a hack/shim method to allow node operators to specify multiple nodes via ENV.
// See: https://app.shortcut.com/chainlinklabs/epic/33587/overhaul-config?cf_workflow=500000005&ct_workflow=all
func SetupNodes(db *sqlx.DB, cfg SetupConfig, lggr logger.Logger) (err error) {
	str := cfg.StarkNetNodes()
	if str == "" {
		return nil
	}

	var nodes []starknetdb.Node
	if err = json.Unmarshal([]byte(str), &nodes); err != nil {
		return errors.Wrapf(err, "invalid STARKNET_NODES json: %q", str)
	}
	// Sorting gives a consistent insert ordering
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

	lggr.Info("STARKNET_NODES was set; clobbering starknet_nodes table")

	orm := NewORM(db, lggr, cfg)
	return orm.SetupNodes(nodes, uniqueIDs(nodes))
}

func uniqueIDs(ns []starknetdb.Node) (ids []string) {
	m := map[string]struct{}{}
	for _, n := range ns {
		id := n.ChainID
		if _, ok := m[id]; ok {
			continue
		}
		ids = append(ids, id)
		m[id] = struct{}{}
	}
	return
}
