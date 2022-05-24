package web

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// ErrTerraNotEnabled is returned when TERRA_ENABLED is not true.
var ErrTerraNotEnabled = errors.New("Terra is disabled. Set TERRA_ENABLED=true to enable.")

func NewTerraNodesController(app chainlink.Application) NodesController {
	parse := func(s string) (string, error) { return s, nil }
	return newNodesController[string, db.Node, presenters.TerraNodeResource](
		app.GetChains().Terra, ErrTerraNotEnabled, parse, presenters.NewTerraNodeResource, func(c *gin.Context) (db.Node, error) {
			var request types.NewNode

			if err := c.ShouldBindJSON(&request); err != nil {
				return db.Node{}, err
			}
			if _, err := app.GetChains().Terra.Show(request.TerraChainID); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					err = fmt.Errorf("Terra chain %s must be added first", request.TerraChainID)
				}
				return db.Node{}, err
			}
			return db.Node{
				Name:          request.Name,
				TerraChainID:  request.TerraChainID,
				TendermintURL: request.TendermintURL,
			}, nil
		})
}
