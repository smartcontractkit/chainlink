package llo

import (
	"context"
	"math/rand/v2"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

var _ LogPoller = (*mockLogPoller)(nil)

type mockLogPoller struct {
	logs        []logpoller.Log
	latestBlock int64

	exprs        []query.Expression
	limitAndSort query.LimitAndSort
}

func (m *mockLogPoller) LatestBlock(ctx context.Context) (logpoller.Block, error) {
	return logpoller.Block{BlockNumber: m.latestBlock}, nil
}
func (m *mockLogPoller) RegisterFilter(ctx context.Context, filter logpoller.Filter) error {
	return nil
}
func (m *mockLogPoller) Replay(ctx context.Context, fromBlock int64) error {
	return nil
}
func (m *mockLogPoller) FilteredLogs(ctx context.Context, filter []query.Expression, limitAndSort query.LimitAndSort, queryName string) ([]logpoller.Log, error) {
	m.exprs = filter
	m.limitAndSort = limitAndSort
	return m.logs, nil
}

type cfg struct {
	cd      ocrtypes.ConfigDigest
	signers [][]byte
	f       uint8
}

type mockConfigCache struct {
	configs map[ocrtypes.ConfigDigest]cfg
}

func (m *mockConfigCache) StoreConfig(ctx context.Context, cd ocrtypes.ConfigDigest, signers [][]byte, f uint8) error {
	m.configs[cd] = cfg{cd, signers, f}
	return nil
}

func Test_ConfigPoller(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger.Test(t)
	lp := &mockLogPoller{make([]logpoller.Log, 0), 0, nil, query.LimitAndSort{}}
	addr := common.Address{1}
	donID := uint32(1)
	donIDHash := DonIDToBytes32(donID)
	fromBlock := uint64(1)

	cc := &mockConfigCache{make(map[ocrtypes.ConfigDigest]cfg)}

	cpBlue := newConfigPoller(lggr, lp, cc, addr, donID, InstanceTypeBlue, fromBlock)
	cpGreen := newConfigPoller(lggr, lp, cc, addr, donID, InstanceTypeGreen, fromBlock)

	cfgCount := uint64(0)

	signers := [][]byte{(common.Address{1}).Bytes(), (common.Address{2}).Bytes()}
	transmitters := []common.Hash{common.Hash{1}, common.Hash{2}}
	f := uint8(1)
	onchainConfig := []byte{5}
	offchainConfigVersion := uint64(6)
	offchainConfig := []byte{7}

	t.Run("Blue/Green config poller follow production and staging configs respectively initially", func(t *testing.T) {
		t.Run("LatestConfigDetails", func(t *testing.T) {
			t.Run("without any logs, returns zero values", func(t *testing.T) {
				changedInBlock, digest, err := cpBlue.LatestConfigDetails(ctx)
				require.NoError(t, err)

				assert.Equal(t, uint64(0), changedInBlock)
				assert.Equal(t, ocr2types.ConfigDigest{}, digest)

				changedInBlock, digest, err = cpGreen.LatestConfigDetails(ctx)
				require.NoError(t, err)

				assert.Equal(t, uint64(0), changedInBlock)
				assert.Equal(t, ocr2types.ConfigDigest{}, digest)
			})

			t.Run("with isGreenProduction=false", func(t *testing.T) {
				isGreenProduction := false

				t.Run("with ProductionConfigSet event, blue starts returning the production config (green still returns zeroes)", func(t *testing.T) {
					cfgCount++
					lp.logs = append(lp.logs, logpoller.Log{Topics: pq.ByteaArray{ProductionConfigSet[:], donIDHash[:]}, Address: addr, EventSig: ProductionConfigSet, BlockNumber: 100, Data: makeConfigSetLogData(
						t,
						0,
						common.Hash{1},
						cfgCount,
						signers,
						transmitters,
						f,
						onchainConfig,
						offchainConfigVersion,
						offchainConfig,
						isGreenProduction,
					)})

					changedInBlock, digest, err := cpBlue.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(100), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{1}, digest)

					changedInBlock, digest, err = cpGreen.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(0), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{}, digest)
				})
				t.Run("with StagingConfigSet event, green starts returning the staging config (blue still returns the production config)", func(t *testing.T) {
					cfgCount++
					lp.logs = append(lp.logs, logpoller.Log{Topics: pq.ByteaArray{StagingConfigSet[:], donIDHash[:]}, Address: addr, EventSig: StagingConfigSet, BlockNumber: 101, Data: makeConfigSetLogData(
						t,
						0,
						common.Hash{2},
						cfgCount,
						signers,
						transmitters,
						f,
						onchainConfig,
						offchainConfigVersion,
						offchainConfig,
						isGreenProduction,
					)})
					changedInBlock, digest, err := cpBlue.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(100), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{1}, digest)

					changedInBlock, digest, err = cpGreen.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(101), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{2}, digest)
				})
			})

			t.Run("with isGreenProduction=true (after PromoteStagingConfig)", func(t *testing.T) {
				isGreenProduction := true
				t.Run("if we ProductionConfigSet again, it now affects green since that is the production instance (blue remains unchanged)", func(t *testing.T) {
					cfgCount++
					lp.logs = append(lp.logs, logpoller.Log{Topics: pq.ByteaArray{ProductionConfigSet[:], donIDHash[:]}, Address: addr, EventSig: ProductionConfigSet, BlockNumber: 103, Data: makeConfigSetLogData(
						t,
						0,
						common.Hash{3},
						cfgCount,
						signers,
						transmitters,
						f,
						onchainConfig,
						offchainConfigVersion,
						offchainConfig,
						isGreenProduction,
					)})

					changedInBlock, digest, err := cpBlue.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(100), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{1}, digest)

					changedInBlock, digest, err = cpGreen.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(103), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{3}, digest)
				})

				t.Run("if we StagingConfigSet again, it now affects blue since that is the staging instance (green remains unchanged)", func(t *testing.T) {
					cfgCount++
					lp.logs = append(lp.logs, logpoller.Log{Topics: pq.ByteaArray{StagingConfigSet[:], donIDHash[:]}, Address: addr, EventSig: StagingConfigSet, BlockNumber: 104, Data: makeConfigSetLogData(
						t,
						0,
						common.Hash{4},
						cfgCount,
						signers,
						transmitters,
						f,
						onchainConfig,
						offchainConfigVersion,
						offchainConfig,
						isGreenProduction,
					)})

					changedInBlock, digest, err := cpBlue.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(104), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{4}, digest)

					changedInBlock, digest, err = cpGreen.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(103), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{3}, digest)
				})

				t.Run("if we StagingConfigSet and ProductionConfigSet again, it sets the config for blue and green respectively", func(t *testing.T) {
					cfgCount++
					lp.logs = append(lp.logs, logpoller.Log{Topics: pq.ByteaArray{StagingConfigSet[:], donIDHash[:]}, Address: addr, EventSig: StagingConfigSet, BlockNumber: 105, Data: makeConfigSetLogData(
						t,
						0,
						common.Hash{5},
						cfgCount,
						signers,
						transmitters,
						f,
						onchainConfig,
						offchainConfigVersion,
						offchainConfig,
						isGreenProduction,
					)})
					cfgCount++
					lp.logs = append(lp.logs, logpoller.Log{Topics: pq.ByteaArray{ProductionConfigSet[:], donIDHash[:]}, Address: addr, EventSig: ProductionConfigSet, BlockNumber: 106, Data: makeConfigSetLogData(
						t,
						0,
						common.Hash{6},
						cfgCount,
						signers,
						transmitters,
						f,
						onchainConfig,
						offchainConfigVersion,
						offchainConfig,
						isGreenProduction,
					)})

					changedInBlock, digest, err := cpBlue.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(105), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{5}, digest)

					changedInBlock, digest, err = cpGreen.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(106), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{6}, digest)
				})
			})
			t.Run("isGreenProduction=false again (another PromoteStagingConfig", func(t *testing.T) {
				isGreenProduction := false
				t.Run("if we PromoteStagingConfig again it re-flips the instances", func(t *testing.T) {
					cfgCount++
					lp.logs = append(lp.logs, logpoller.Log{LogIndex: 1, Topics: pq.ByteaArray{ProductionConfigSet[:], donIDHash[:]}, Address: addr, EventSig: ProductionConfigSet, BlockNumber: 107, Data: makeConfigSetLogData(
						t,
						0,
						common.Hash{7},
						cfgCount,
						signers,
						transmitters,
						f,
						onchainConfig,
						offchainConfigVersion,
						offchainConfig,
						isGreenProduction,
					)})
					cfgCount++
					lp.logs = append(lp.logs, logpoller.Log{LogIndex: 2, Topics: pq.ByteaArray{StagingConfigSet[:], donIDHash[:]}, Address: addr, EventSig: StagingConfigSet, BlockNumber: 107, Data: makeConfigSetLogData(
						t,
						0,
						common.Hash{8},
						cfgCount,
						signers,
						transmitters,
						f,
						onchainConfig,
						offchainConfigVersion,
						offchainConfig,
						isGreenProduction,
					)})

					changedInBlock, digest, err := cpBlue.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(107), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{7}, digest)

					changedInBlock, digest, err = cpGreen.LatestConfigDetails(ctx)
					require.NoError(t, err)

					assert.Equal(t, uint64(107), changedInBlock)
					assert.Equal(t, ocr2types.ConfigDigest{8}, digest)
				})
			})
			t.Run("Stores all seen configs in cache", func(t *testing.T) {
				require.Len(t, cc.configs, 8)
				for i := 1; i <= 8; i++ {
					assert.Contains(t, cc.configs, ocr2types.ConfigDigest{byte(i)})
					assert.Equal(t, cfg{
						cd:      ocr2types.ConfigDigest{byte(i)},
						signers: signers,
						f:       f,
					}, cc.configs[ocr2types.ConfigDigest{byte(i)}])
				}
			})
		})
		t.Run("LatestConfig", func(t *testing.T) {
			t.Run("changedInBlock corresponds to a block in which a log was emitted, returns the config", func(t *testing.T) {
				expectedSigners := []ocr2types.OnchainPublicKey{ocr2types.OnchainPublicKey{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, ocr2types.OnchainPublicKey{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}
				expectedTransmitters := []ocr2types.Account{"0100000000000000000000000000000000000000000000000000000000000000", "0200000000000000000000000000000000000000000000000000000000000000"}

				cfg, err := cpBlue.LatestConfig(ctx, 107)
				require.NoError(t, err)
				assert.Equal(t, ocr2types.ConfigDigest{7}, cfg.ConfigDigest)
				assert.Equal(t, uint64(7), cfg.ConfigCount)
				assert.Equal(t, expectedSigners, cfg.Signers)
				assert.Equal(t, expectedTransmitters, cfg.Transmitters)
				assert.Equal(t, f, cfg.F)
				assert.Equal(t, onchainConfig, cfg.OnchainConfig)
				assert.Equal(t, offchainConfigVersion, cfg.OffchainConfigVersion)
				assert.Equal(t, offchainConfig, cfg.OffchainConfig)

				cfg, err = cpGreen.LatestConfig(ctx, 107)
				require.NoError(t, err)
				assert.Equal(t, ocr2types.ConfigDigest{8}, cfg.ConfigDigest)
				assert.Equal(t, uint64(8), cfg.ConfigCount)
				assert.Equal(t, expectedSigners, cfg.Signers)
				assert.Equal(t, expectedTransmitters, cfg.Transmitters)
				assert.Equal(t, f, cfg.F)
				assert.Equal(t, onchainConfig, cfg.OnchainConfig)
				assert.Equal(t, offchainConfigVersion, cfg.OffchainConfigVersion)
				assert.Equal(t, offchainConfig, cfg.OffchainConfig)
			})
		})
		t.Run("LatestBlockHeight", func(t *testing.T) {
			t.Run("returns the latest block from log poller", func(t *testing.T) {
				latest := rand.Int64()
				lp.latestBlock = latest

				latestBlock, err := cpBlue.LatestBlockHeight(ctx)
				require.NoError(t, err)
				assert.Equal(t, uint64(latest), latestBlock)

				latestBlock, err = cpGreen.LatestBlockHeight(ctx)
				require.NoError(t, err)
				assert.Equal(t, uint64(latest), latestBlock)
			})
		})
	})
}

func makeConfigSetLogData(t *testing.T,
	previousConfigBlockNumber uint32,
	configDigest common.Hash,
	configCount uint64,
	signers [][]byte,
	offchainTransmitters []common.Hash,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	isGreenProduction bool,
) []byte {
	event := configuratorABI.Events["ProductionConfigSet"]
	data, err := event.Inputs.NonIndexed().Pack(previousConfigBlockNumber, configDigest, configCount, signers, offchainTransmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, isGreenProduction)
	require.NoError(t, err)
	return data
}
