package llo

import (
	"bytes"
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

type ShouldRetireCacheService interface {
	services.Service
	llotypes.ShouldRetireCache
}

type shouldRetireCache struct {
	services.Service
	eng *services.Engine

	lp          LogPoller
	addr        common.Address
	donID       uint32
	donIDTopic  common.Hash
	filterExprs []query.Expression

	pollPeriod time.Duration

	mu             sync.RWMutex
	latestBlockNum int64
	m              map[ocrtypes.ConfigDigest]struct{}
}

func NewShouldRetireCache(lggr logger.Logger, lp LogPoller, addr common.Address, donID uint32) ShouldRetireCacheService {
	return newShouldRetireCache(lggr, lp, addr, donID)
}

func newShouldRetireCache(lggr logger.Logger, lp LogPoller, addr common.Address, donID uint32) *shouldRetireCache {
	donIDTopic := DonIDToBytes32(donID)
	exprs := []query.Expression{
		logpoller.NewAddressFilter(addr),
		logpoller.NewEventSigFilter(PromoteStagingConfig),
		logpoller.NewEventByTopicFilter(1, []logpoller.HashedValueComparator{
			{Values: []common.Hash{donIDTopic}, Operator: primitives.Eq},
		}),
		// NOTE: Optimize for fast retirement detection. On Arbitrum,
		// finalization can take tens of minutes
		// (https://grafana.ops.prod.cldev.sh/d/e0453cc9-4b4a-41e1-9f01-7c21de805b39/blockchain-finality-and-gas?orgId=1&var-env=All&var-network_name=ethereum-testnet-sepolia-arbitrum-1&var-network_name=ethereum-mainnet-arbitrum-1&from=1732460992641&to=1732547392641)
		query.Confidence(primitives.Unconfirmed),
	}
	s := &shouldRetireCache{
		lp:          lp,
		addr:        addr,
		donID:       donID,
		donIDTopic:  donIDTopic,
		filterExprs: exprs,
		m:           make(map[ocrtypes.ConfigDigest]struct{}),
		pollPeriod:  1 * time.Second,
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

	exprs := make([]query.Expression, 0, len(s.filterExprs)+1)
	exprs = append(exprs, s.filterExprs...)
	exprs = append(exprs,
		query.Block(strconv.FormatInt(fromBlock, 10), primitives.Gte),
	)

	logs, err := s.lp.FilteredLogs(ctx, exprs, NoLimitSortAsc, "ShouldRetireCache - PromoteStagingConfig")
	if err != nil {
		s.eng.SugaredLogger.Errorw("checkShouldRetire: IndexedLogs", "err", err)
		return
	}

	for _, log := range logs {
		if log.EventSig != PromoteStagingConfig {
			// ignore unrecognized logs
			continue
		}

		if !bytes.Equal(log.Topics[1], s.donIDTopic[:]) {
			// skip logs for other donIDs, shouldn't happen given the
			// FilterLogs call, but belts and braces
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
