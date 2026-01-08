package devenv

import (
	"context"
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/devenv/products/ocr2"
)

type Cfg struct {
	ProductType string              `toml:"product_type"`
	Blockchains []*blockchain.Input `toml:"blockchains" validate:"required"`
	FakeServer  *fake.Input         `toml:"fake_server" validate:"required"`
	NodeSets    []*ns.Input         `toml:"nodesets"    validate:"required"`
	JD          *jd.Input           `toml:"jd"`
}

func newProduct(typ string) (Product, error) {
	switch typ {
	case "ocr2":
		return ocr2.NewOCR2Configurator(), nil
	default:
		return nil, fmt.Errorf("unknown product type: %s", typ)
	}
}

func NewEnvironment(ctx context.Context) error {
	if err := framework.DefaultNetwork(nil); err != nil {
		return err
	}
	in, err := Load[Cfg]()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	_, err = blockchain.NewBlockchainNetwork(in.Blockchains[0])
	if err != nil {
		return fmt.Errorf("failed to create blockchain network 1337: %w", err)
	}
	if os.Getenv("FAKE_SERVER_IMAGE") != "" {
		in.FakeServer.Image = os.Getenv("FAKE_SERVER_IMAGE")
	}
	_, err = fake.NewDockerFakeDataProvider(in.FakeServer)
	if err != nil {
		return fmt.Errorf("failed to create fake data provider: %w", err)
	}

	c, err := newProduct(in.ProductType)
	if err != nil {
		return err
	}
	if err = c.Load(); err != nil {
		return fmt.Errorf("failed to load product config: %w", err)
	}

	overrides, err := c.GenerateCLNodesBlockchainConfig(ctx, in.Blockchains[0])
	if err != nil {
		return fmt.Errorf("failed to generate CL nodes config: %w", err)
	}
	for _, ns := range in.NodeSets[0].NodeSpecs {
		ns.Node.TestConfigOverrides = overrides
		if os.Getenv("CHAINLINK_IMAGE") != "" {
			ns.Node.Image = os.Getenv("CHAINLINK_IMAGE")
		}
	}

	_, err = ns.NewSharedDBNodeSet(in.NodeSets[0], nil)
	if err != nil {
		return fmt.Errorf("failed to create new shared db node set: %w", err)
	}

	err = c.ConfigureJobsAndContracts(
		ctx,
		in.FakeServer,
		in.Blockchains[0],
		in.NodeSets[0],
	)
	if err != nil {
		return fmt.Errorf("failed to setup default product deployment: %w", err)
	}
	L.Info().Str("BootstrapNode", in.NodeSets[0].Out.CLNodes[0].Node.ExternalURL).Send()
	for _, n := range in.NodeSets[0].Out.CLNodes[1:] {
		L.Info().Str("Node", n.Node.ExternalURL).Send()
	}
	if err := Store[Cfg](in); err != nil {
		return fmt.Errorf("failed to write infra config: %w", err)
	}
	return c.Store("env-out.toml")
}
