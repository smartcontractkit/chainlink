package arbiter

import (
	"context"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	// ShardConfigContractName is the name used to identify the ShardConfig contract.
	ShardConfigContractName = "ShardConfig"

	// GetDesiredShardCountMethod is the method name for reading the desired shard count.
	GetDesiredShardCountMethod = "getDesiredShardCount"
)

// ShardConfigABI is the ABI for the ShardConfig contract.
// Only includes the getDesiredShardCount function we need.
const ShardConfigABI = `[
	{
		"inputs": [],
		"name": "getDesiredShardCount",
		"outputs": [{"internalType": "uint256", "name": "", "type": "uint256"}],
		"stateMutability": "view",
		"type": "function"
	}
]`

// ShardConfigReader reads the desired shard count from the ShardConfig contract.
type ShardConfigReader interface {
	// GetDesiredShardCount retrieves the current desired shard count from on-chain.
	GetDesiredShardCount(ctx context.Context) (uint64, error)
}

type shardConfigReader struct {
	reader   types.ContractReader
	contract types.BoundContract
	lggr     logger.Logger
}

// NewShardConfigReader creates a new ShardConfigReader.
func NewShardConfigReader(
	reader types.ContractReader,
	contractAddress string,
	lggr logger.Logger,
) ShardConfigReader {
	return &shardConfigReader{
		reader: reader,
		contract: types.BoundContract{
			Address: contractAddress,
			Name:    ShardConfigContractName,
		},
		lggr: lggr.Named("ShardConfigReader"),
	}
}

// GetDesiredShardCount retrieves the current desired shard count from on-chain.
func (s *shardConfigReader) GetDesiredShardCount(ctx context.Context) (uint64, error) {
	var result *big.Int

	err := s.reader.GetLatestValue(
		ctx,
		s.contract.ReadIdentifier(GetDesiredShardCountMethod),
		primitives.Finalized, // Use finalized for governance data
		nil,                  // No input params
		&result,
	)
	if err != nil {
		s.lggr.Errorw("failed to get desired shard count from on-chain",
			"error", err,
			"contract", s.contract.Address,
		)
		return 0, err
	}

	count := result.Uint64()

	// Update metrics
	SetOnChainMaxReplicas(count)

	s.lggr.Debugw("read desired shard count from on-chain",
		"count", count,
		"contract", s.contract.Address,
	)

	return count, nil
}
