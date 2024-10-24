package llo

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func Test_OffchainConfigDigester_ConfigDigest(t *testing.T) {
	ctx := tests.Context(t)
	// ChainID and ContractAddress are taken into account for computation
	cd1, err := OffchainConfigDigester{ChainID: big.NewInt(0)}.ConfigDigest(ctx, types.ContractConfig{})
	require.NoError(t, err)
	cd2, err := OffchainConfigDigester{ChainID: big.NewInt(0)}.ConfigDigest(ctx, types.ContractConfig{})
	require.NoError(t, err)
	cd3, err := OffchainConfigDigester{ChainID: big.NewInt(1)}.ConfigDigest(ctx, types.ContractConfig{})
	require.NoError(t, err)
	cd4, err := OffchainConfigDigester{ChainID: big.NewInt(1), ContractAddress: common.Address{1}}.ConfigDigest(ctx, types.ContractConfig{})
	require.NoError(t, err)

	require.Equal(t, cd1, cd2)
	require.NotEqual(t, cd2, cd3)
	require.NotEqual(t, cd2, cd4)
	require.NotEqual(t, cd3, cd4)

	configID := common.HexToHash("0x1")
	chainID := big.NewInt(2)
	addr := common.HexToAddress("0x3")
	prefix := ocrtypes.ConfigDigestPrefix(4)

	digester := NewOffchainConfigDigester(configID, chainID, addr, prefix)
	// any signers ok
	_, err = digester.ConfigDigest(ctx, types.ContractConfig{
		Signers: []types.OnchainPublicKey{{1, 2}},
	})
	require.NoError(t, err)

	// malformed transmitters
	_, err = digester.ConfigDigest(ctx, types.ContractConfig{
		Transmitters: []types.Account{"0x"},
	})
	require.Error(t, err)

	_, err = digester.ConfigDigest(ctx, types.ContractConfig{
		Transmitters: []types.Account{"7343581f55146951b0f678dc6cfa8fd360e2f353"},
	})
	require.Error(t, err)

	_, err = digester.ConfigDigest(ctx, types.ContractConfig{
		Transmitters: []types.Account{"7343581f55146951b0f678dc6cfa8fd360e2f353aabbccddeeffaaccddeeffaz"},
	})
	require.Error(t, err)

	// well-formed transmitters
	_, err = digester.ConfigDigest(ctx, types.ContractConfig{
		Transmitters: []types.Account{"7343581f55146951b0f678dc6cfa8fd360e2f353aabbccddeeffaaccddeeffaa"},
	})
	require.NoError(t, err)
}
