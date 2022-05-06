package web

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/web/presenters"
)

// ErrSolanaNotEnabled is returned when SOLANA_ENABLED is not true.
var ErrSolanaNotEnabled = errors.New("Solana is disabled. Set SOLANA_ENABLED=true to enable.")

func NewSolanaNodesController(app chainlink.Application) NodesController {
	parse := func(s string) (string, error) { return s, nil }
	return newNodesController[string, db.Node, presenters.SolanaNodeResource](
		app.GetChains().Solana, ErrSolanaNotEnabled, parse, presenters.NewSolanaNodeResource, func(c *gin.Context) (db.Node, error) {
			var request db.NewNode

			if err := c.ShouldBindJSON(&request); err != nil {
				return db.Node{}, err
			}
			if _, err := app.GetChains().Solana.Show(request.SolanaChainID); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					err = errors.Errorf("Solana chain %s must be added first", request.SolanaChainID)
				}
				return db.Node{}, err
			}
			return db.Node{
				Name:          request.Name,
				SolanaChainID: request.SolanaChainID,
				SolanaURL:     request.SolanaURL,
			}, nil
		})
}
