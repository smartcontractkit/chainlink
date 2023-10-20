package test_env

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/logwatch"
	"github.com/smartcontractkit/chainlink-testing-framework/networks"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"

	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/types/config/node"
)

type CleanUpType string

const (
	CleanUpTypeNone     CleanUpType = "none"
	CleanUpTypeStandard CleanUpType = "standard"
	CleanUpTypeCustom   CleanUpType = "custom"
)

type CLTestEnvBuilder struct {
	hasLogWatch        bool
	hasGeth            bool
	hasKillgrave       bool
	hasForwarders      bool
	clNodeConfig       *chainlink.Config
	secretsConfig      string
	nonDevGethNetworks []blockchain.EVMNetwork
	clNodesCount       int
	customNodeCsaKeys  []string
	defaultNodeCsaKeys []string
	l                  zerolog.Logger
	t                  *testing.T
	te                 *CLClusterTestEnv
	isNonEVM           bool
	cleanUpType        CleanUpType
	cleanUpCustomFn    func()

	/* funding */
	ETHFunds *big.Float
}

func NewCLTestEnvBuilder() *CLTestEnvBuilder {
	return &CLTestEnvBuilder{
		l: log.Logger,
	}
}

// WithTestEnv sets the test environment to use for the test.
// If nil, a new test environment is created.
// If not nil, the test environment is used as-is.
// If TEST_ENV_CONFIG_PATH is set, the test environment is created with the config at that path.
func (b *CLTestEnvBuilder) WithTestEnv(te *CLClusterTestEnv) (*CLTestEnvBuilder, error) {
	envConfigPath, isSet := os.LookupEnv("TEST_ENV_CONFIG_PATH")
	var cfg *TestEnvConfig
	var err error
	if isSet {
		cfg, err = NewTestEnvConfigFromFile(envConfigPath)
		if err != nil {
			return nil, err
		}
	}

	if te != nil {
		b.te = te
	} else {
		b.te, err = NewTestEnv()
		if err != nil {
			return nil, err
		}
	}

	if cfg != nil {
		b.te = b.te.WithTestEnvConfig(cfg)
	}
	return b, nil
}

// WithTestLogger sets the test logger to use for the test.
// Useful for parallel tests so the logging will be separated correctly in the results views.
func (b *CLTestEnvBuilder) WithTestLogger(t *testing.T) *CLTestEnvBuilder {
	b.t = t
	b.l = logging.GetTestLogger(t)
	return b
}

func (b *CLTestEnvBuilder) WithLogWatcher() *CLTestEnvBuilder {
	b.hasLogWatch = true
	return b
}

func (b *CLTestEnvBuilder) WithCLNodes(clNodesCount int) *CLTestEnvBuilder {
	b.clNodesCount = clNodesCount
	return b
}

func (b *CLTestEnvBuilder) WithForwarders() *CLTestEnvBuilder {
	b.hasForwarders = true
	return b
}

func (b *CLTestEnvBuilder) WithFunding(eth *big.Float) *CLTestEnvBuilder {
	b.ETHFunds = eth
	return b
}

func (b *CLTestEnvBuilder) WithGeth() *CLTestEnvBuilder {
	b.hasGeth = true
	return b
}

func (b *CLTestEnvBuilder) WithPrivateGethChains(evmNetworks []blockchain.EVMNetwork) *CLTestEnvBuilder {
	b.nonDevGethNetworks = evmNetworks
	return b
}

func (b *CLTestEnvBuilder) WithCLNodeConfig(cfg *chainlink.Config) *CLTestEnvBuilder {
	b.clNodeConfig = cfg
	return b
}

func (b *CLTestEnvBuilder) WithSecretsConfig(secrets string) *CLTestEnvBuilder {
	b.secretsConfig = secrets
	return b
}

func (b *CLTestEnvBuilder) WithMockAdapter() *CLTestEnvBuilder {
	b.hasKillgrave = true
	return b
}

// WithNonEVM sets the test environment to not use EVM when built.
func (b *CLTestEnvBuilder) WithNonEVM() *CLTestEnvBuilder {
	b.isNonEVM = true
	return b
}

func (b *CLTestEnvBuilder) WithStandardCleanup() *CLTestEnvBuilder {
	b.cleanUpType = CleanUpTypeStandard
	return b
}

func (b *CLTestEnvBuilder) WithoutCleanup() *CLTestEnvBuilder {
	b.cleanUpType = CleanUpTypeNone
	return b
}

func (b *CLTestEnvBuilder) WithCustomCleanup(customFn func()) *CLTestEnvBuilder {
	b.cleanUpType = CleanUpTypeCustom
	b.cleanUpCustomFn = customFn
	return b
}

func (b *CLTestEnvBuilder) Build() (*CLClusterTestEnv, error) {
	if b.te == nil {
		var err error
		b, err = b.WithTestEnv(nil)
		if err != nil {
			return nil, err
		}
	}
	b.l.Info().
		Bool("hasGeth", b.hasGeth).
		Bool("hasKillgrave", b.hasKillgrave).
		Int("clNodesCount", b.clNodesCount).
		Strs("customNodeCsaKeys", b.customNodeCsaKeys).
		Strs("defaultNodeCsaKeys", b.defaultNodeCsaKeys).
		Msg("Building CL cluster test environment..")

	var err error
	if b.t != nil {
		b.te.WithTestLogger(b.t)
	}

	if b.hasLogWatch {
		b.te.LogWatch, err = logwatch.NewLogWatch(nil, nil)
		if err != nil {
			return nil, err
		}
	}

	if b.hasKillgrave {
		err = b.te.StartMockAdapter()
		if err != nil {
			return nil, err
		}
	}

	switch b.cleanUpType {
	case CleanUpTypeStandard:
		b.t.Cleanup(func() {
			if err := b.te.Cleanup(); err != nil {
				b.l.Error().Err(err).Msg("Error cleaning up test environment")
			}
		})
	case CleanUpTypeCustom:
		b.t.Cleanup(b.cleanUpCustomFn)
	case CleanUpTypeNone:
		b.l.Warn().Msg("test environment won't be cleaned up")
	case "":
		return b.te, errors.WithMessage(errors.New("explicit cleanup type must be set when building test environment"), "test environment builder failed")
	}

	if b.nonDevGethNetworks != nil {
		b.te.WithPrivateChain(b.nonDevGethNetworks)
		err := b.te.StartPrivateChain()
		if err != nil {
			return b.te, err
		}
		var nonDevNetworks []blockchain.EVMNetwork
		for i, n := range b.te.PrivateChain {
			primaryNode := n.GetPrimaryNode()
			if primaryNode == nil {
				return b.te, errors.WithStack(fmt.Errorf("primary node is nil in PrivateChain interface"))
			}
			nonDevNetworks = append(nonDevNetworks, *n.GetNetworkConfig())
			nonDevNetworks[i].URLs = []string{primaryNode.GetInternalWsUrl()}
			nonDevNetworks[i].HTTPURLs = []string{primaryNode.GetInternalHttpUrl()}
		}
		if nonDevNetworks == nil {
			return nil, errors.New("cannot create nodes with custom config without nonDevNetworks")
		}

		err = b.te.StartClCluster(b.clNodeConfig, b.clNodesCount, b.secretsConfig)
		if err != nil {
			return nil, err
		}
		return b.te, nil
	}
	networkConfig := networks.SelectedNetwork
	var internalDockerUrls test_env.InternalDockerUrls
	if b.hasGeth && networkConfig.Simulated {
		networkConfig, internalDockerUrls, err = b.te.StartGeth()
		if err != nil {
			return nil, err
		}

	}

	if !b.isNonEVM {
		bc, err := blockchain.NewEVMClientFromNetwork(networkConfig, b.l)
		if err != nil {
			return nil, err
		}

		b.te.EVMClient = bc
		cd, err := contracts.NewContractDeployer(bc, b.l)
		if err != nil {
			return nil, err
		}
		b.te.ContractDeployer = cd

		cl, err := contracts.NewContractLoader(bc, b.l)
		if err != nil {
			return nil, err
		}
		b.te.ContractLoader = cl
	}

	var nodeCsaKeys []string

	// Start Chainlink Nodes
	if b.clNodesCount > 0 {
		var cfg *chainlink.Config
		if b.clNodeConfig != nil {
			cfg = b.clNodeConfig
		} else {
			cfg = node.NewConfig(node.NewBaseConfig(),
				node.WithOCR1(),
				node.WithP2Pv1(),
			)
		}

		if !b.isNonEVM {
			var httpUrls []string
			var wsUrls []string
			if networkConfig.Simulated {
				httpUrls = []string{internalDockerUrls.HttpUrl}
				wsUrls = []string{internalDockerUrls.WsUrl}
			} else {
				httpUrls = networkConfig.HTTPURLs
				wsUrls = networkConfig.URLs
			}

			node.SetChainConfig(cfg, wsUrls, httpUrls, networkConfig, b.hasForwarders)
		}

		err := b.te.StartClCluster(cfg, b.clNodesCount, b.secretsConfig)
		if err != nil {
			return nil, err
		}

		nodeCsaKeys, err = b.te.ClCluster.NodeCSAKeys()
		if err != nil {
			return nil, err
		}
		b.defaultNodeCsaKeys = nodeCsaKeys
	}

	if b.hasGeth && b.clNodesCount > 0 && b.ETHFunds != nil {
		b.te.ParallelTransactions(true)
		defer b.te.ParallelTransactions(false)
		if err := b.te.FundChainlinkNodes(b.ETHFunds); err != nil {
			return nil, err
		}
	}

	return b.te, nil
}
