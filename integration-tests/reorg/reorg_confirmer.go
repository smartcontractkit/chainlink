package reorg

import (
	"context"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/chaos"
	"github.com/smartcontractkit/chainlink-env/environment"
	a "github.com/smartcontractkit/chainlink-env/pkg/alias"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/reorg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
)

// The steps are:
// 1. Await initial network consensus for N blocks
// 2. Fork the network, separate node 0 from 1 and 2
// 3. Await alternative blocks mined on node zero,
// nodes 1 and 2 will have longer chain with more block difficulty so blocks of node 0 must be replaced
// 4. Await consensus for N blocks, check alt block hashes now
const (
	InitConsensus int64 = iota
	NetworkFork
	CheckBlocks
	Consensus
	Wait
)

type ReorgConfig struct {
	FromPodLabel            string
	ToPodLabel              string
	Network                 blockchain.EVMClient
	Env                     *environment.Environment
	BlockConsensusThreshold int
	Timeout                 time.Duration
}

// ReorgController reorg stats collecting struct
type ReorgController struct {
	cfg                   *ReorgConfig
	ReorgDepth            int
	numberOfNodes         int
	currentBlockConsensus int
	forkBlockNumber       int64
	currentAltBlocks      int
	currentVerifiedBlocks int
	networkStep           atomic.Int64
	altBlockNumbers       []int64
	blocksByNode          map[int]map[int64]blockchain.NodeHeader
	blockHashes           map[int64][]common.Hash
	chaosExperimentName   string
	initConsensusReady    chan struct{}
	reorgStarted          chan struct{}
	depthReached          chan struct{}
	once                  *sync.Once
	mutex                 sync.Mutex
	ctx                   context.Context
	cancel                context.CancelFunc
	doneChan              chan struct{}
	complete              bool
}

// NewReorgController creates a type that can create reorg chaos and confirm reorg has happened
func NewReorgController(cfg *ReorgConfig) (*ReorgController, error) {
	if len(cfg.Network.GetClients()) == 1 {
		return nil, errors.New("need at least 3 nodes to re-org")
	}
	ctx, ctxCancel := context.WithTimeout(context.Background(), cfg.Timeout)
	rc := &ReorgController{
		cfg:                cfg,
		numberOfNodes:      len(cfg.Network.GetClients()),
		doneChan:           make(chan struct{}, 100),
		altBlockNumbers:    make([]int64, 0),
		blockHashes:        map[int64][]common.Hash{},
		blocksByNode:       map[int]map[int64]blockchain.NodeHeader{},
		initConsensusReady: make(chan struct{}, 100),
		reorgStarted:       make(chan struct{}, 100),
		depthReached:       make(chan struct{}, 100),
		once:               &sync.Once{},
		mutex:              sync.Mutex{},
		ctx:                ctx,
		cancel:             ctxCancel,
		complete:           false,
	}
	rc.networkStep.Store(InitConsensus)
	cfg.Network.AddHeaderEventSubscription("reorg", rc)
	<-rc.initConsensusReady
	return rc, nil
}

// ReOrg start re-org
func (rc *ReorgController) ReOrg(depth int) {
	rc.ReorgDepth = depth
	rc.networkStep.Store(NetworkFork)
}

// WaitDepthReached wait until N alternative blocks are mined
func (rc *ReorgController) WaitDepthReached() error {
	<-rc.depthReached
	rc.networkStep.Store(Consensus)
	return rc.cfg.Network.WaitForEvents()
}

// WaitReorgStarted waits until first alternative block is present
func (rc *ReorgController) WaitReorgStarted() {
	<-rc.reorgStarted
}

// ReceiveHeader receives header marked by node that mined it,
// forks the network and record all alternative block numbers on node 0
func (rc *ReorgController) ReceiveHeader(header blockchain.NodeHeader) error {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.appendBlockHeader(header)
	switch rc.networkStep.Load() {
	case Wait:
	case InitConsensus:
		if rc.hasNetworkFormedConsensus(header) {
			rc.networkStep.Store(Wait)
			rc.initConsensusReady <- struct{}{}
		}
	case NetworkFork:
		if err := rc.forkNetwork(header); err != nil {
			return err
		}
	case CheckBlocks:
		if err := rc.compareBlocks(header); err != nil {
			return err
		}
	case Consensus:
		if rc.hasNetworkFormedConsensus(header) {
			rc.complete = true
			rc.cancel()
		}
	}
	return nil
}

func (rc *ReorgController) Complete() bool {
	return rc.complete
}

// VerifyReorgComplete verifies that all blocks are replaced by reorg
func (rc *ReorgController) VerifyReorgComplete() error {
	// we are excluding the last block because we can't control joining while blocks were mined
	for _, ab := range rc.altBlockNumbers[:len(rc.altBlockNumbers)-1] {
		bb := rc.blocksByNode[0][ab]
		h, err := rc.cfg.Network.HeaderHashByNumber(context.Background(), big.NewInt(ab))
		if err != nil {
			return err
		}
		log.Info().
			Int64("Number", bb.Number.Int64()).
			Str("Hash before", bb.Hash().String()).
			Str("Hash after", h).
			Msg("Comparing block")
		if bb.Hash().String() != h {
			rc.currentVerifiedBlocks++
		}
	}
	if rc.currentVerifiedBlocks+1 < rc.ReorgDepth {
		return errors.New("Reorg depth has not met")
	}
	return nil
}

func (rc *ReorgController) isAltBlock(blk blockchain.NodeHeader) bool {
	blockNumber := blk.Number.Int64()
	// If we've received the same block number from all nodes, check hashes there are some different versions
	if len(rc.blockHashes[blockNumber]) >= rc.numberOfNodes {
		firstBlockHash := rc.blockHashes[blockNumber][0]
		for _, blockHash := range rc.blockHashes[blockNumber][1:] {
			if blockHash.String() != firstBlockHash.String() {
				return true
			}
		}
	}
	return false
}

func (rc *ReorgController) compareBlocks(blk blockchain.NodeHeader) error {
	if blk.NodeID == 0 && blk.Number.Int64() >= rc.forkBlockNumber && rc.isAltBlock(blk) {
		rc.once.Do(func() {
			log.Warn().Int64("Number", blk.Number.Int64()).Msg("Reorg started")
			rc.reorgStarted <- struct{}{}
		})
		rc.altBlockNumbers = append(rc.altBlockNumbers, blk.Number.Int64())
		rc.currentAltBlocks++
		log.Info().
			Int64("Number", blk.Number.Int64()).
			Str("Hash", blk.Hash().String()).
			Int("Node", blk.NodeID).
			Int("BlocksLeft", rc.ReorgDepth-rc.currentAltBlocks).
			Msg("Mined alternative block")
	}
	if rc.currentAltBlocks == rc.ReorgDepth {
		log.Info().Msg("Joining network")
		if err := rc.joinNetwork(); err != nil {
			return err
		}
		rc.depthReached <- struct{}{}
		rc.currentBlockConsensus = 0
	}
	return nil
}

// Wait wait until reorg is done
func (rc *ReorgController) Wait() error {
	<-rc.ctx.Done()
	rc.cfg.Network.DeleteHeaderEventSubscription("reorg")
	if rc.complete {
		return nil
	}
	return errors.New("timeout waiting for reorg to complete")
}

// forkNetwork stomp the network between target reorged node and the rest
func (rc *ReorgController) forkNetwork(header blockchain.NodeHeader) error {
	rc.forkBlockNumber = header.Number.Int64()
	log.Debug().
		Int64("Number", rc.forkBlockNumber).
		Str("Network", rc.cfg.Network.GetNetworkName()).
		Msg("Forking network")
	expName, err := rc.cfg.Env.Chaos.Run(
		chaos.NewNetworkPartition(
			rc.cfg.Env.Cfg.Namespace,
			&chaos.Props{
				DurationStr: "999h",
				FromLabels:  &map[string]*string{"app": a.Str(reorg.TXNodesAppLabel)},
				ToLabels:    &map[string]*string{"app": a.Str(reorg.MinerNodesAppLabel)},
			},
		))
	rc.chaosExperimentName = expName
	rc.networkStep.Store(CheckBlocks)
	return err
}

// joinNetwork restores network connectivity between nodes
func (rc *ReorgController) joinNetwork() error {
	return rc.cfg.Env.Chaos.Stop(rc.chaosExperimentName)
}

// hasNetworkFormedConsensus waits for network to have N blocks with the same hashes
func (rc *ReorgController) hasNetworkFormedConsensus(header blockchain.NodeHeader) bool {
	blockNumber := header.Number.Int64()
	// If we've received the same block number from all nodes, check hashes to ensure they've reformed consensus
	if len(rc.blockHashes[blockNumber]) >= rc.numberOfNodes {
		firstBlockHash := rc.blockHashes[blockNumber][0]
		for _, blockHash := range rc.blockHashes[blockNumber][1:] {
			if blockHash.String() != firstBlockHash.String() {
				return false
			}
		}
		log.Debug().
			Int("Blocks left", rc.cfg.BlockConsensusThreshold-rc.currentBlockConsensus).
			Msg("Waiting for consensus threshold")
		rc.currentBlockConsensus++
	}

	if rc.currentBlockConsensus >= rc.cfg.BlockConsensusThreshold {
		log.Info().
			Msg("Network joined")
		return true
	}
	return false
}

func (rc *ReorgController) appendBlockHeader(header blockchain.NodeHeader) {
	bn := header.Number.Int64()
	if _, ok := rc.blockHashes[bn]; !ok {
		rc.blockHashes[bn] = []common.Hash{}
	}
	rc.blockHashes[bn] = append(rc.blockHashes[bn], header.Hash())

	if _, ok := rc.blocksByNode[header.NodeID]; !ok {
		rc.blocksByNode[header.NodeID] = make(map[int64]blockchain.NodeHeader)
	}
	rc.blocksByNode[header.NodeID][bn] = header
}
