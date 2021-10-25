package chaos

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/chaos/experiments"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/environment"
)

// ReorgStep step enum for re-org test
type ReorgStep int

// The steps are:
// 1. Await initial network consensus for N blocks
// 2. Fork the network, separate node 0 from 1 and 2
// 3. Await alternative blocks mined on node zero,
// nodes 1 and 2 will have longer chain with more block difficulty so blocks of node 0 must be replaced
// 4. Await consensus for N blocks, check alt block hashes now
const (
	InitConsensus ReorgStep = iota
	NetworkFork
	AltBlocks
	Consensus
)

// Must use that labels for all networks, TXNode is the main node connected to Chainlink
const (
	TXNodeLabel    = "tx-node"
	MinerNodeLabel = "miner-node"
)

// ReorgConfirmer reorg stats collecting struct
type ReorgConfirmer struct {
	env environment.Environment
	c   client.BlockchainClient
	// number of blocks we want to be reorged on node zero (default client)
	reorgDepth int
	// blocks for that all nodes sent equal hashes, amount of blocks to wait
	blockConsensusThreshold int
	numberOfNodes           int
	currentBlockConsensus   int
	forkBlockNumber         int64
	currentAltBlocks        int
	currentVerifiedBlocks   int
	NetworkStep             ReorgStep
	altBlockNumbers         []int64
	blocksByNode            map[int]map[int64]client.NodeBlock
	blockHashes             map[int64][]common.Hash
	chaosExperimentName     string
	mutex                   sync.Mutex
	ctx                     context.Context
	cancel                  context.CancelFunc
	doneChan                chan struct{}
	done                    bool
}

// NewReorgConfirmer creates a type that can create reorg chaos and confirm reorg has happened
func NewReorgConfirmer(
	c client.BlockchainClient,
	env environment.Environment,
	reorgDepth int,
	blockConsensusThreshold int,
	timeout time.Duration,
) (*ReorgConfirmer, error) {
	if len(c.GetClients()) == 1 {
		return nil, errors.New("only one node within the blockchain client detected, cannot reorg")
	}
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	rc := &ReorgConfirmer{
		env:                     env,
		c:                       c,
		reorgDepth:              reorgDepth,
		blockConsensusThreshold: blockConsensusThreshold,
		numberOfNodes:           len(c.GetClients()),
		doneChan:                make(chan struct{}),
		altBlockNumbers:         make([]int64, 0),
		blockHashes:             map[int64][]common.Hash{},
		blocksByNode:            map[int]map[int64]client.NodeBlock{},
		mutex:                   sync.Mutex{},
		ctx:                     ctx,
		cancel:                  ctxCancel,
	}
	return rc, nil
}

// ReceiveBlock receives block marked by node that mined it,
// forks the network and record all alternative block numbers on node 0
func (rc *ReorgConfirmer) ReceiveBlock(blk client.NodeBlock) error {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	if blk.Block == nil || rc.done {
		return nil
	}
	rc.appendBlockHeader(blk)
	switch rc.NetworkStep {
	case InitConsensus:
		if rc.hasNetworkFormedConsensus(blk) {
			rc.NetworkStep = NetworkFork
		}
	case NetworkFork:
		if err := rc.forkNetwork(blk); err != nil {
			return err
		}
	case AltBlocks:
		if err := rc.awaitAltBlocks(blk); err != nil {
			return err
		}
	case Consensus:
		if rc.hasNetworkFormedConsensus(blk) {
			rc.done = true
			rc.doneChan <- struct{}{}
		}
	}
	return nil
}

// Verify verifies that all blocks are replaced by reorg
func (rc *ReorgConfirmer) Verify() error {
	// we are excluding the last block because we can't control joining while blocks were mined
	for _, ab := range rc.altBlockNumbers[:len(rc.altBlockNumbers)-1] {
		bb := rc.blocksByNode[0][ab]
		h, err := rc.c.HeaderHashByNumber(context.Background(), big.NewInt(ab))
		if err != nil {
			return err
		}
		log.Info().
			Int64("Number", bb.Number().Int64()).
			Str("Hash before", bb.Hash().String()).
			Str("Hash after", h).
			Msg("Comparing block")
		if bb.Hash().String() != h {
			rc.currentVerifiedBlocks++
		}
	}
	if rc.currentVerifiedBlocks+1 < rc.reorgDepth {
		return errors.New("Reorg depth has not met")
	}
	return nil
}

func (rc *ReorgConfirmer) isAltBlock(blk client.NodeBlock) bool {
	blockNumber := blk.Number().Int64()
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

func (rc *ReorgConfirmer) awaitAltBlocks(blk client.NodeBlock) error {
	if blk.NodeID == 0 && blk.Number().Int64() >= rc.forkBlockNumber && rc.isAltBlock(blk) {
		log.Info().
			Int64("Number", blk.Number().Int64()).
			Str("Hash", blk.Hash().String()).
			Int("Node", blk.NodeID).
			Msg("Mined alternative block")
		rc.altBlockNumbers = append(rc.altBlockNumbers, blk.Number().Int64())
		rc.currentAltBlocks++
	}
	if rc.currentAltBlocks == rc.reorgDepth {
		log.Info().Msg("Joining network")
		if err := rc.joinNetwork(); err != nil {
			return err
		}
		rc.NetworkStep = Consensus
		rc.currentBlockConsensus = 0
	}
	return nil
}

// Wait wait until reorg is done
func (rc *ReorgConfirmer) Wait() error {
	for {
		select {
		case <-rc.doneChan:
			rc.cancel()
			return nil
		case <-rc.ctx.Done():
			if rc.done {
				return nil
			}
			return errors.New("timeout waiting for reorg to complete")
		}
	}
}

func (rc *ReorgConfirmer) forkNetwork(blk client.NodeBlock) error {
	rc.forkBlockNumber = blk.Number().Int64()
	log.Debug().Int64("Number", rc.forkBlockNumber).Msg("Forking network")
	expName, err := rc.env.ApplyChaos(&experiments.NetworkPartition{
		FromMode:       "one",
		FromLabelKey:   "app",
		FromLabelValue: TXNodeLabel,
		ToMode:         "all",
		ToLabelKey:     "app",
		ToLabelValue:   MinerNodeLabel,
	})
	rc.chaosExperimentName = expName
	rc.NetworkStep = AltBlocks
	return err
}

func (rc *ReorgConfirmer) joinNetwork() error {
	return rc.env.StopChaos(rc.chaosExperimentName)
}

func (rc *ReorgConfirmer) hasNetworkFormedConsensus(blk client.NodeBlock) bool {
	blockNumber := blk.Number().Int64()
	// If we've received the same block number from all nodes, check hashes to ensure they've reformed consensus
	if len(rc.blockHashes[blockNumber]) >= rc.numberOfNodes {
		firstBlockHash := rc.blockHashes[blockNumber][0]
		for _, blockHash := range rc.blockHashes[blockNumber][1:] {
			if blockHash.String() != firstBlockHash.String() {
				return false
			}
		}
		log.Debug().
			Int("Blocks left", rc.blockConsensusThreshold-rc.currentBlockConsensus).
			Msg("Waiting for consensus threshold")
		rc.currentBlockConsensus++
	}

	if rc.currentBlockConsensus >= rc.blockConsensusThreshold && !rc.done {
		log.Info().
			Msg("Network joined")
		return true
	}
	return false
}

func (rc *ReorgConfirmer) appendBlockHeader(blk client.NodeBlock) {
	bn := blk.Number().Int64()
	if _, ok := rc.blockHashes[bn]; !ok {
		rc.blockHashes[bn] = []common.Hash{}
	}
	rc.blockHashes[bn] = append(rc.blockHashes[bn], blk.Hash())

	if _, ok := rc.blocksByNode[blk.NodeID]; !ok {
		rc.blocksByNode[blk.NodeID] = make(map[int64]client.NodeBlock)
	}
	rc.blocksByNode[blk.NodeID][bn] = blk
}
