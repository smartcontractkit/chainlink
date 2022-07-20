package web

import (
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewTerraChainsController(app chainlink.Application) ChainsController {
	parse := func(s string) (string, error) { return s, nil }
	return newChainsController[string, *db.ChainCfg, presenters.TerraChainResource](
		"terra", app.GetChains().Terra, ErrTerraNotEnabled, parse, presenters.NewTerraChainResource)
}
