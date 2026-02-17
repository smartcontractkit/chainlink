package devenv

import (
	"context"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"

	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

// Product describes a minimal set of methods that each legacy product must implement
type Product interface {
	// Load describes how to load your product-specific config
	Load() error

	// Store describes how to store your product-specific config output
	// The output may include URLs to you services, CLDF contracts addresses and more
	Store(path string, instanceIdx int) error

	// GenerateNodesSecrets describes how to generate secrets for Chainlink nodes for your product
	// deployed on multiple blockchains and nodesets
	GenerateNodesSecrets(
		ctx context.Context,
		fs *fake.Input,
		bc []*blockchain.Input,
		ns []*nodeset.Input,
	) (string, error)

	// GenerateNodesConfig describes how to generate Chainlink node config
	// specific to your product deployed on multiple blockchains and nodesets
	GenerateNodesConfig(
		ctx context.Context,
		fs *fake.Input,
		bc []*blockchain.Input,
		ns []*nodeset.Input,
	) (string, error)

	// ConfigureJobsAndContracts describe how to configure jobs and contracts
	// specifically to your product deployed on multiple blockchains and nodesets
	// Configuration may be called multiple times if "instances" key is specified in "env.toml"
	// the implementation should be aware of it and be able to configure multiple instances of the product
	ConfigureJobsAndContracts(
		ctx context.Context,
		instanceIdx int,
		fs *fake.Input,
		bc []*blockchain.Input,
		ns []*nodeset.Input,
	) error
}
