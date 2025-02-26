package deployment_ethereum

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestOpGasBump(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironmentWithOps(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain1 := e.AllChainSelectors()[0]

	client := e.Chains[chain1].Client
	auth := e.Chains[chain1].DeployerKey

	deps := EthereumDeps{
		Auth:    auth,
		Client:  client,
		Confirm: e.Chains[chain1].ConfirmByHash,
	}

	input := GasBumpOpInput{
		GasBump{
			RetryLimit:      10,
			RetryIntervalMs: 1000,
			BumpPercentage:  10,
		},
		struct {
			To    *common.Address
			Data  []byte
			Value *big.Int
		}{
			To:    &auth.From,
			Value: big.NewInt(1000000000000000),
		},
	}

	report, err := deployment.ExecuteOp(e.OpEnv, SendTxWithGasBumpOp, deps, input)
	require.NoError(t, err)

	// pretty print the report with json intentation
	reportJSON, _ := json.MarshalIndent(report, "", "  ")
	t.Log(string(reportJSON))
}
