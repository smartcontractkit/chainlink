package devenv

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/devenv/products/ocr2"
)

type ProductInfo struct {
	Name      string `toml:"name"`
	Instances int    `toml:"instances"`
}

type Cfg struct {
	Products    []*ProductInfo      `toml:"products"`
	Blockchains []*blockchain.Input `toml:"blockchains" validate:"required"`
	FakeServer  *fake.Input         `toml:"fake_server" validate:"required"`
	NodeSets    []*ns.Input         `toml:"nodesets"    validate:"required"`
	JD          *jd.Input           `toml:"jd"`
}

func newProduct(name string) (Product, error) {
	switch name {
	case "ocr2":
		return ocr2.NewOCR2Configurator(), nil
	default:
		return nil, fmt.Errorf("unknown product type: %s", name)
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

	// get all the product orchestrations, generate product specific overrides
	productConfigurators := make([]Product, 0)
	clNodeProductOverrides := make([]string, 0)
	for _, product := range in.Products {
		p, err := newProduct(product.Name)
		if err != nil {
			return err
		}
		if err = p.Load(); err != nil {
			return fmt.Errorf("failed to load product config: %w", err)
		}

		overrides, err := p.GenerateCLNodesBlockchainConfig(ctx, in.Blockchains[0])
		if err != nil {
			return fmt.Errorf("failed to generate CL nodes config: %w", err)
		}

		productConfigurators = append(productConfigurators, p)
		clNodeProductOverrides = append(clNodeProductOverrides, overrides)
	}

	// merge overrides, spin up node sets and write infrastructure outputs
	// infra is always common for all the products, if it can't be we should fail
	// user should use different infra layout in env.toml then
	for _, ns := range in.NodeSets[0].NodeSpecs {
		ns.Node.TestConfigOverrides = strings.Join(clNodeProductOverrides, "\n")
		if os.Getenv("CHAINLINK_IMAGE") != "" {
			ns.Node.Image = os.Getenv("CHAINLINK_IMAGE")
		}
	}
	_, err = ns.NewSharedDBNodeSet(in.NodeSets[0], nil)
	if err != nil {
		return fmt.Errorf("failed to create new shared db node set: %w", err)
	}
	if err := Store[Cfg](in); err != nil {
		return err
	}

	// deploy all products and all instances,
	// product config function controls what to read and how to orchestrate each instance
	// via their own TOML part, we only deploy N instances of product M
	for productIdx, productInfo := range in.Products {
		for productInstance := range productInfo.Instances {
			err = productConfigurators[productIdx].ConfigureJobsAndContracts(
				ctx,
				in.FakeServer,
				in.Blockchains[0],
				in.NodeSets[0],
			)
			if err != nil {
				return fmt.Errorf("failed to setup default product deployment: %w", err)
			}
			if err := productConfigurators[productIdx].Store("env-out.toml", productInstance); err != nil {
				return errors.New("failed to store product config")
			}
		}
	}
	L.Info().Str("BootstrapNode", in.NodeSets[0].Out.CLNodes[0].Node.ExternalURL).Send()
	for _, n := range in.NodeSets[0].Out.CLNodes[1:] {
		L.Info().Str("Node", n.Node.ExternalURL).Send()
	}
	return nil
}
