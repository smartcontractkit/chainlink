package deployment_common

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLinkChangeset(t *testing.T) {

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain1 := e.AllChainSelectors()[0]

	changesetInput := ChangesetLinkInput{
		MintAmount: big.NewInt(1000000000000000000),
		Amount:     big.NewInt(1000000000000),
		To:         common.HexToAddress("0x1"),
		chainID:    chain1,
	}
	ret, err := LinkExampleChangeset(e, changesetInput)
	require.NoError(t, err)

	jsonReports, err := json.MarshalIndent(ret.Reports, "", "  ")
	if err != nil {
		t.Errorf("Failed to marshal reports to JSON: %v", err)
	}
	t.Log(string(jsonReports))

	require.NoError(t, err)
	require.NotNil(t, ret)
}
