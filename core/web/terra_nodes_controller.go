package web

import (
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// ErrTerraNotEnabled is returned when TERRA_ENABLED is not true.
var ErrTerraNotEnabled = errChainDisabled{name: "Terra", envVar: "TERRA_ENABLED"}

func NewTerraNodesController(app chainlink.Application) NodesController {
	parse := func(s string) (string, error) { return s, nil }
	return newNodesController[string, db.Node, presenters.TerraNodeResource](
		app.GetChains().Terra, ErrTerraNotEnabled, parse, presenters.NewTerraNodeResource, app.GetAuditLogger())
}
