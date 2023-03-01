package web

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink/core/chains/cosmos/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// ErrCosmosNotEnabled is returned when COSMOS_ENABLED is not true.
var ErrCosmosNotEnabled = errChainDisabled{name: "Cosmos", envVar: "COSMOS_ENABLED"}

func NewCosmosNodesController(app chainlink.Application) NodesController {
	parse := func(s string) (string, error) { return s, nil }
	return newNodesController[string, db.Node, presenters.CosmosNodeResource](
		app.GetChains().Cosmos, ErrCosmosNotEnabled, parse, presenters.NewCosmosNodeResource, func(c *gin.Context) (db.Node, error) {
			var request types.NewNode

			if err := c.ShouldBindJSON(&request); err != nil {
				return db.Node{}, err
			}
			if _, err := app.GetChains().Cosmos.Show(request.CosmosChainID); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					err = fmt.Errorf("Cosmos chain %s must be added first", request.CosmosChainID)
				}
				return db.Node{}, err
			}
			return db.Node{
				Name:          request.Name,
				CosmosChainID: request.CosmosChainID,
				TendermintURL: request.TendermintURL,
			}, nil
		},
		app.GetLogger(),
		app.GetAuditLogger(),
	)
}
