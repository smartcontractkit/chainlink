package devenv

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink/devenv/products/ocr2"
)

type Cfg struct {
	OCR2        *ocr2.OCR2          `toml:"ocr2"`
	Blockchains []*blockchain.Input `toml:"blockchains" validate:"required"`
	FakeServer  *fake.Input         `toml:"fake_server" validate:"required"`
	NodeSets    []*ns.Input         `toml:"nodesets"    validate:"required"`
	JD          *jd.Input           `toml:"jd"`
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
	_, err = fake.NewDockerFakeDataProvider(in.FakeServer)
	if err != nil {
		return fmt.Errorf("failed to create fake data provider: %w", err)
	}

	pkey := getNetworkPrivateKey()
	if pkey == "" {
		return fmt.Errorf("PRIVATE_KEY environment variable not set")
	}
	bc := in.Blockchains[0].Out.Nodes[0]
	in.OCR2.OCR2DynamicConfig = &ocr2.OCR2DynamicConfig{
		PKeyStr:                   pkey,
		ChainID:                   in.Blockchains[0].ChainID,
		FakeServerExternalHTTPURL: in.FakeServer.Out.BaseURLHost,
		FakeServerInternalHTTPURL: in.FakeServer.Out.BaseURLDocker,
		BlockchainExternalWSURL:   bc.ExternalWSUrl,
		BlockchainInternalWSURL:   bc.InternalWSUrl,
		BlockchainInternalHTTPURL: bc.InternalHTTPUrl,
	}

	overrides, err := ocr2.ConfigureCLNodes(ctx, in.OCR2)
	if err != nil {
		return fmt.Errorf("failed to setup default CLDF orchestration: %w", err)
	}
	for _, ns := range in.NodeSets[0].NodeSpecs {
		ns.Node.TestConfigOverrides = overrides
	}
	_, err = ns.NewSharedDBNodeSet(in.NodeSets[0], nil)
	if err != nil {
		return fmt.Errorf("failed to create new shared db node set: %w", err)
	}
	in.OCR2.OCR2DynamicConfig.BootstrapContainerName = in.NodeSets[0].Out.CLNodes[0].Node.ContainerName

	clClients, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	if err != nil {
		return err
	}
	_, err = ocr2.ConfigureContractsAndJobs(
		ctx,
		clClients,
		in.OCR2,
		ocr2.ConfigureProductContractsJobs,
	)
	if err != nil {
		return fmt.Errorf("failed to setup default CLDF orchestration: %w", err)
	}
	L.Info().Str("BootstrapNode", in.NodeSets[0].Out.CLNodes[0].Node.ExternalURL).Send()
	for _, n := range in.NodeSets[0].Out.CLNodes[1:] {
		L.Info().Str("Node", n.Node.ExternalURL).Send()
	}
	return Store[Cfg](in)
}
