package job

import (
	"testing"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func NewTestORM(t *testing.T, db *sqlx.DB, chainSet evm.ChainSet, pipelineORM pipeline.ORM, keyStore keystore.Master) ORM {
	o := NewORM(db, chainSet, pipelineORM, keyStore, logger.TestLogger(t))
	t.Cleanup(func() { o.Close() })
	return o
}
