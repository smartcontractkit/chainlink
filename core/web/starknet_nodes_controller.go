package web

import (
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

// ErrStarknetNotEnabled is returned when Starknet.Enabled is not true.
var ErrStarknetNotEnabled = errChainDisabled{name: "Starknet", envVar: "Starknet.Enabled"}

func NewStarknetNodesController(app chainlink.Application) NodesController {
	parse := func(s string) (string, error) { return s, nil }
	return newNodesController[string, db.Node, presenters.StarknetNodeResource](
		app.GetChains().Starknet, ErrStarknetNotEnabled, parse, presenters.NewStarknetNodeResource, app.GetAuditLogger())
}
