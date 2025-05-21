package example

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployAndMintExampleChangeset(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain1 := e.AllChainSelectors()[0]

	changesetInput := SqDeployLinkInput{
		MintAmount: big.NewInt(1000000000000000000),
		Amount:     big.NewInt(1000000000000),
		To:         common.HexToAddress("0x1"),
		ChainID:    chain1,
	}
	result, err := DeployAndMintExampleChangeset{}.Apply(e, changesetInput)
	require.NoError(t, err)

	require.Len(t, result.Reports, 4) // 3 ops + 1 seq report
	require.NoError(t, err)
}
