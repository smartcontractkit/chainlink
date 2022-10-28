package web

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// ErrStarkNetNotEnabled is returned when STARKNET_ENABLED is not true.
var ErrStarkNetNotEnabled = errChainDisabled{name: "StarkNet", envVar: "STARKNET_ENABLED"}

func NewStarkNetNodesController(app chainlink.Application) NodesController {
	parse := func(s string) (string, error) { return s, nil }
	return newNodesController[string, db.Node, presenters.StarkNetNodeResource](
		app.GetChains().StarkNet, ErrStarkNetNotEnabled, parse, presenters.NewStarkNetNodeResource, func(c *gin.Context) (db.Node, error) {
			var request struct {
				Name    string
				ChainID string
				URL     string
			}

			if err := c.ShouldBindJSON(&request); err != nil {
				return db.Node{}, err
			}
			if _, err := app.GetChains().StarkNet.Show(request.ChainID); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					err = errors.Errorf("StarkNet chain %s must be added first", request.ChainID)
				}
				return db.Node{}, err
			}
			return db.Node{
				Name:    request.Name,
				ChainID: request.ChainID,
				URL:     request.URL,
			}, nil
		}, app.GetLogger(), app.GetAuditLogger())
}
