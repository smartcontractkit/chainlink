package example

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployAndMintExampleChangeset(t *testing.T) {
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain1 := e.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))[0]

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
