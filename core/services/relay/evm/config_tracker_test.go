package evm_test

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
)

func TestConfigTracker_LatestConfig(t *testing.T) {
	lggr := logger.TestLogger(t)
	c := new(evmmocks.Client)
	contractABI, _ := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	ct := evm.NewConfigTracker(lggr, contractABI, c, common.Address{}, "", nil)

	configSet, _ := hex.DecodeString(
		"0000000000000000000000000000000000000000000000000000000000000000000168dbbc989af81ad798fdc0102dcd7608a0f3943f5372d70e066f4cc47aef0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001c000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000260000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000002c00000000000000000000000000000000000000000000000000000000000000004000000000000000000000000100a209e68e25ecf836a1563b34b0bd884de3c9c000000000000000000000000c9e3fe3a9b0ab463abc67f510001b52b3ee2fabc000000000000000000000000f9cbaa5e7d007a828c1d903ec827fa73c5813cf30000000000000000000000005dc70bd55d9a5300d8dfd2b50307ecf62eaef33600000000000000000000000000000000000000000000000000000000000000040000000000000000000000008ce423450190e1069f146b0d36ccef5a111c34ca00000000000000000000000027f93e255054e5435ed66ce29bba464f49104c970000000000000000000000000ada788ab72408cd016baed9966929157fe0888b0000000000000000000000004d6524c5afe74d68e0bc5c6edaf74a08423e13a90000000000000000000000000000000000000000000000000000000000000031018000000000000000000000000000000000000000000000007fffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000023b0880a8d6b907108094ebdc03188094ebdc032080cab5ee012880a8d6b90730033a04010101014220741efdafc6e2f6ea87b58e33619b15cd0f0fb8432c7c5b17bebe64a4d2744648422092806cf7f5ce0c8ac298bd8d168a91424601d1d1de9b36bff6cf3c44d159f4ef422035674e56b0e46e4bcab213c46ef1ba95a6cd9278440bf35ae5e358afd32c605742203108fb8086b6338e92519b7cd70846db44ad35613096feda392c661083a6708c4a34313244334b6f6f574a375a3176394469733351765261476d5a33666e4b79375646724b64717474585a5a56735a356d534b58716e4a34313244334b6f6f574753545944486f5a595148646759614437523161396232484e69466f4448656f4854706e4d6e6b7a783733444a34313244334b6f6f574c3971724150575253596e55757444587658527555316269585842424776315a7145726f6f737442354a39334a34313244334b6f6f574465437471455144436154744171697a47765367444d726f3979345574674d526d595659596f466641776939520a1080ade2042080ade2045880e1eb176080e1eb176880e1eb177080e1eb177880e1eb1782018c010a20befed1f74bcd9dc9342cd6106d78604a214d9904a37ed0b6058595f605e9783012201a1bf753c2505dab5204b284f1fefd7dc4eb83f6304664c9c2f43cbf425666601a10b16b9ed7f188a4381d81a6ae8d379abd1a10d6a941106211fcaf568ccc1493a0a2161a105a765837db073bd319cbd6e6cadfa9d91a1015cecfc264d0ad9567f7b5c3e44e0f950000000000")
	c.On("FilterLogs", mock.Anything, mock.Anything).Return([]types.Log{{Topics: []common.Hash{contractABI.Events["ConfigSet"].ID}, Data: configSet}}, nil).Once()
	cfg, err := ct.LatestConfig(context.Background(), 10)
	require.NoError(t, err)
	// Spot check a few values
	assert.Equal(t, uint8(1), cfg.F)
	assert.Equal(t, 4, len(cfg.Signers))
	assert.Equal(t, 4, len(cfg.Transmitters))
}

func Test_OCRContractTracker_LatestBlockHeight(t *testing.T) {
	t.Parallel()

	t.Run("on L2 chains, always returns 0", func(t *testing.T) {
		uni := newContractTrackerUni(t, evmtest.ChainOptimismMainnet(t))
		l, err := uni.configTracker.LatestBlockHeight(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(0), l)
	})

	t.Run("before first head incoming, looks up on-chain", func(t *testing.T) {
		uni := newContractTrackerUni(t)
		uni.ec.On("HeadByNumber", mock.AnythingOfType("*context.cancelCtx"), (*big.Int)(nil)).Return(&evmtypes.Head{Number: 42}, nil)

		l, err := uni.configTracker.LatestBlockHeight(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(42), l)
	})

	t.Run("Before first head incoming, on client error returns error", func(t *testing.T) {
		uni := newContractTrackerUni(t)
		uni.ec.On("HeadByNumber", mock.AnythingOfType("*context.cancelCtx"), (*big.Int)(nil)).Return(nil, nil).Once()

		_, err := uni.configTracker.LatestBlockHeight(context.Background())
		assert.EqualError(t, err, "got nil head")

		uni.ec.On("HeadByNumber", mock.AnythingOfType("*context.cancelCtx"), (*big.Int)(nil)).Return(nil, errors.New("bar")).Once()

		_, err = uni.configTracker.LatestBlockHeight(context.Background())
		assert.EqualError(t, err, "bar")

		uni.ec.AssertExpectations(t)
	})

	t.Run("after first head incoming, uses cached value", func(t *testing.T) {
		uni := newContractTrackerUni(t)

		uni.configTracker.OnNewLongestChain(context.Background(), &evmtypes.Head{Number: 42})

		l, err := uni.configTracker.LatestBlockHeight(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(42), l)
	})

	t.Run("if headbroadcaster has it, uses the given value on start", func(t *testing.T) {
		uni := newContractTrackerUni(t)

		uni.hb.On("Subscribe", uni.configTracker).Return(&evmtypes.Head{Number: 42}, func() {})
		require.NoError(t, uni.configTracker.Start())

		l, err := uni.configTracker.LatestBlockHeight(context.Background())
		require.NoError(t, err)

		assert.Equal(t, uint64(42), l)

		uni.hb.AssertExpectations(t)

		require.NoError(t, uni.configTracker.Close())
	})
}
