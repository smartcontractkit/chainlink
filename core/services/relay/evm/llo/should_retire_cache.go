package llo

import (
	"bytes"
	"context"
	"math"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

type ShouldRetireCacheService interface {
	services.Service
	llotypes.ShouldRetireCache
}

type shouldRetireCache struct {
	services.Service
	eng *services.Engine

	lp        LogPoller
	addr      common.Address
	donID     uint32
	donIDHash common.Hash

	pollPeriod time.Duration

	mu             sync.RWMutex
	latestBlockNum int64
	m              map[ocrtypes.ConfigDigest]struct{}
}

func NewShouldRetireCache(lggr logger.Logger, lp LogPoller, addr common.Address, donID uint32) ShouldRetireCacheService {
	return newShouldRetireCache(lggr, lp, addr, donID)
}

func newShouldRetireCache(lggr logger.Logger, lp LogPoller, addr common.Address, donID uint32) *shouldRetireCache {
	s := &shouldRetireCache{
		lp:         lp,
		addr:       addr,
		donID:      donID,
		donIDHash:  DonIDToBytes32(donID),
		m:          make(map[ocrtypes.ConfigDigest]struct{}),
		pollPeriod: 1 * time.Second,
	}
	s.Service, s.eng = services.Config{
		Name:  "LLOShouldRetireCache",
		Start: s.start,
	}.NewServiceEngine(lggr)

	return s
}

func (s *shouldRetireCache) ShouldRetire(digest ocr2types.ConfigDigest) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.m[digest]
	if exists {
		s.eng.SugaredLogger.Debugw("ShouldRetire", "digest", digest, "shouldRetire", exists)
	}
	return exists, nil
}

func (s *shouldRetireCache) start(ctx context.Context) error {
	t := services.TickerConfig{
		Initial:   0,
		JitterPct: services.DefaultJitter,
	}.NewTicker(s.pollPeriod)
	s.eng.GoTick(t, s.checkShouldRetire)
	return nil
}

func (s *shouldRetireCache) checkShouldRetire(ctx context.Context) {
	fromBlock := s.latestBlockNum + 1
	logs, err := s.lp.LogsWithSigs(ctx, fromBlock, math.MaxInt64, []common.Hash{PromoteStagingConfig}, s.addr)
	if err != nil {
		s.eng.SugaredLogger.Errorw("checkShouldRetire: IndexedLogs", "err", err)
		return
	}

	for _, log := range logs {
		// TODO: This can probably be optimized
		// MERC-3524
		if !bytes.Equal(log.Topics[1], s.donIDHash[:]) {
			continue
		}
		digestBytes := log.Topics[2]
		digest, err := ocrtypes.BytesToConfigDigest(digestBytes)
		if err != nil {
			s.eng.SugaredLogger.Errorw("checkShouldRetire: BytesToConfigDigest failed", "err", err)
			return
		}
		s.eng.SugaredLogger.Infow("markRetired: Got retired config digest", "blockNum", log.BlockNumber, "configDigest", digest)
		s.markRetired(digest)
		if log.BlockNumber > s.latestBlockNum {
			s.latestBlockNum = log.BlockNumber
		}
	}
}

func (s *shouldRetireCache) markRetired(configDigest ocrtypes.ConfigDigest) {
	s.mu.RLock()
	_, exists := s.m[configDigest]
	s.mu.RUnlock()
	if exists {
		// already marked as retired
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[configDigest] = struct{}{}
}
