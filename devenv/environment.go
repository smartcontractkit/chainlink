package devenv

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"

	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

type Cfg struct {
	OCR2        *OCR2               `toml:"ocr2"`
	Blockchains []*blockchain.Input `toml:"blockchains" validate:"required"`
	FakeServer  *fake.Input         `toml:"fake_server" validate:"required"`
	NodeSets    []*ns.Input         `toml:"nodesets"    validate:"required"`
	JD          *jd.Input           `toml:"jd"`
}

func NewEnvironment() (*Cfg, error) {
	if err := framework.DefaultNetwork(nil); err != nil {
		return nil, err
	}
	in, err := Load[Cfg]()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	_, err = blockchain.NewBlockchainNetwork(in.Blockchains[0])
	if err != nil {
		return nil, fmt.Errorf("failed to create blockchain network 1337: %w", err)
	}
	_, err = fake.NewDockerFakeDataProvider(in.FakeServer)
	if err != nil {
		return nil, fmt.Errorf("failed to create fake data provider: %w", err)
	}
	if err := OCR2ProductConfiguration(in, ConfigureNodesNetwork); err != nil {
		return nil, fmt.Errorf("failed to setup default CLDF orchestration: %w", err)
	}
	_, err = ns.NewSharedDBNodeSet(in.NodeSets[0], nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new shared db node set: %w", err)
	}
	if err := OCR2ProductConfiguration(in, ConfigureProductContractsJobs); err != nil {
		return nil, fmt.Errorf("failed to setup default CLDF orchestration: %w", err)
	}
	return in, Store[Cfg](in)
}
