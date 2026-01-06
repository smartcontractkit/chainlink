package devenv

import (
	"context"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"

	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

// Product describes a minimal set of methods that each legacy product must implement
type Product interface {
	// Load loads product-specific config part from TOML
	Load() error
	// Store stores product-specific config part to TOML
	Store(path string) error
	// GenerateCLNodesBlockchainConfig generates configuration for CL nodes for blockchain connection
	GenerateCLNodesBlockchainConfig(
		ctx context.Context,
		bc *blockchain.Input,
	) (string, error)
	// ConfigureJobsAndContracts configures both on-chain and off-chain parts of a product
	ConfigureJobsAndContracts(
		ctx context.Context,
		fs *fake.Input,
		bc *blockchain.Input,
		ns *nodeset.Input,
	) error
}
