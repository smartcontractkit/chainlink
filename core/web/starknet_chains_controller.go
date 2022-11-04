package web

import (
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

func NewStarkNetChainsController(app chainlink.Application) ChainsController {
	return newChainsController[string, *db.ChainCfg]("starknet", app.GetChains().StarkNet, ErrStarkNetNotEnabled,
		func(s string) (string, error) { return s, nil }, presenters.NewStarkNetChainResource, app.GetLogger(), app.GetAuditLogger())
}
