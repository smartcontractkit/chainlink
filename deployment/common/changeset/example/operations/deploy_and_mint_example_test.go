package example

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
)

func TestDeployAndMintExampleChangeset(t *testing.T) {
	selector := chain_selectors.TEST_90000001.Selector
	env, err := environment.New(t.Context(), environment.WithEVMSimulated(t, []uint64{selector}))
	require.NoError(t, err)

	changesetInput := SqDeployLinkInput{
		MintAmount: big.NewInt(1000000000000000000),
		Amount:     big.NewInt(1000000000000),
		To:         common.HexToAddress("0x1"),
		ChainID:    selector,
	}
	result, err := DeployAndMintExampleChangeset{}.Apply(*env, changesetInput)
	require.NoError(t, err)

	require.Len(t, result.Reports, 4) // 3 ops + 1 seq report
	require.NoError(t, err)
}
