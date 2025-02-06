package deployment_common

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	deployment_ethereum "github.com/smartcontractkit/chainlink/deployment/ethereum/extension"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestLinkOps(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain1 := e.AllChainSelectors()[0]

	client := e.Chains[chain1].Client
	auth := e.Chains[chain1].DeployerKey

	ctx := deployment.OpContext{
		Log: lggr,
	}

	deps := deployment_ethereum.EthereumDeps{
		Auth:    auth,
		Client:  client,
		Confirm: e.Chains[chain1].ConfirmByHash,
	}

	// DEPLOY
	deployRes, err := DeployLinkOp.Execute(ctx, deps, deployment.EmptyInput{})
	require.NoError(t, err)
	require.NotNil(t, deployRes.Address)
	require.NotNil(t, deployRes.Hash)

	// GRANT MINT ROLE
	grantMint, err := GrantMintLinkOp.Execute(ctx, deps, GrantLinkInput{
		contractAddress: deployRes.Address,
		To:              auth.From,
	})
	require.NoError(t, err)
	require.NotNil(t, grantMint.Hash)

	// MINT SOME TO SELF
	mint, err := MintLinkOp.Execute(ctx, deps, MintLinkInput{
		contractAddress: deployRes.Address,
		To:              auth.From,
		Amount:          big.NewInt(1000000000000000000),
	})
	require.NoError(t, err)
	require.NotNil(t, mint.Hash)

	// TRANSFER SOME TO DEST
	dest := common.HexToAddress("0x1")
	transferInput := TransferLinkInput{
		contractAddress: deployRes.Address,
		To:              dest,
		Amount:          big.NewInt(1000000000000000),
	}

	transferRes, err := TransferLinkOp.Execute(ctx, deps, transferInput)
	require.NoError(t, err)
	require.NotNil(t, transferRes.Hash)

	link, err := link_token.NewLinkToken(deployRes.Address, client)

	// CHECK BALANCE
	balanceOfDest, err := link.BalanceOf(nil, dest)
	require.NoError(t, err)
	require.NotNil(t, balanceOfDest)

}
